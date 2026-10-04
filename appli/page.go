package groupie

import (
	"fmt"
	"image/color"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var maFenetre fyne.Window
var tousLesArtistes []List_artist
var lApplication fyne.App

var filtresActuels FilterOptions
var rechercheActuelle string

type FilterOptions struct {
	CreationMin   int
	CreationMax   int
	FirstAlbumMin int
	FirstAlbumMax int
	MembersCounts []int
	Locations     []string
}

func OuvertureApp() {
	lApplication = app.New()
	lApplication.Settings().SetTheme(theme.DarkTheme())

	maFenetre = lApplication.NewWindow("Groupie Tracker nael-joey-zayed")
	maFenetre.Resize(fyne.NewSize(1200, 800))
	maFenetre.CenterOnScreen()

	dictionnaireArtistes := Api()
	for _, artiste := range dictionnaireArtistes {
		tousLesArtistes = append(tousLesArtistes, artiste)
	}
	sort.Slice(tousLesArtistes, func(i, j int) bool {
		return tousLesArtistes[i].Name < tousLesArtistes[j].Name
	})

	AfficherLeMenu()
	maFenetre.ShowAndRun()
}

func AfficherLeMenu() {
	if desk, ok := maFenetre.Canvas().(desktop.Canvas); ok {
		desk.SetOnKeyDown(nil)
	}

	entryRecherche := widget.NewEntry()
	entryRecherche.SetPlaceHolder("Rechercher un artiste, membre, année, lieu...")
	entryRecherche.SetText(rechercheActuelle)

	grilleArtistes := container.NewGridWithColumns(4)
	conteneurDefilant := container.NewVScroll(grilleArtistes)

	rafraichirGrille := func() {
		grilleArtistes.Objects = nil

		artistesAffiches := FiltrerArtistes(tousLesArtistes, filtresActuels)
		artistesAffiches = RechercherArtistes(artistesAffiches, rechercheActuelle)

		for _, a := range artistesAffiches {
			var img *canvas.Image
			if uri, err := storage.ParseURI(a.Image); err == nil {
				img = canvas.NewImageFromURI(uri)
			} else {
				img = canvas.NewImageFromResource(theme.MediaMusicIcon())
			}
			img.FillMode = canvas.ImageFillContain
			img.SetMinSize(fyne.NewSize(150, 150))

			imgCliq := NewImageCliquable(img, func(art List_artist) func() {
				return func() { AfficherLesDetails(art) }
			}(a))

			lbl := widget.NewLabel(a.Name)
			lbl.Alignment = fyne.TextAlignCenter
			lbl.TextStyle = fyne.TextStyle{Bold: true}

			grilleArtistes.Add(container.NewVBox(imgCliq, lbl))
		}
		grilleArtistes.Refresh()
	}

	entryRecherche.OnChanged = func(texte string) {
		rechercheActuelle = texte
		rafraichirGrille()
	}

	barreRecherche := container.NewBorder(nil, nil,
		widget.NewIcon(theme.SearchIcon()),
		nil,
		entryRecherche,
	)

	entryCreationMin := widget.NewEntry()
	entryCreationMin.SetPlaceHolder("Année création min")
	entryCreationMax := widget.NewEntry()
	entryCreationMax.SetPlaceHolder("Année création max")

	entryAlbumMin := widget.NewEntry()
	entryAlbumMin.SetPlaceHolder("Premier album min (YYYY)")
	entryAlbumMax := widget.NewEntry()
	entryAlbumMax.SetPlaceHolder("Premier album max (YYYY)")

	chk1 := widget.NewCheck("1 membre", nil)
	chk2 := widget.NewCheck("2 membres", nil)
	chk3 := widget.NewCheck("3 membres", nil)
	chk4 := widget.NewCheck("4+ membres", nil)

	entryLocations := widget.NewEntry()
	entryLocations.SetPlaceHolder("Lieux concerts (séparés par , ex: Paris, USA)")

	btnAppliquer := widget.NewButton("Appliquer les filtres", func() {
		parseVal := func(s string) int {
			if v, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				return v
			}
			return 0
		}

		filtresActuels = FilterOptions{
			CreationMin:   parseVal(entryCreationMin.Text),
			CreationMax:   parseVal(entryCreationMax.Text),
			FirstAlbumMin: parseVal(entryAlbumMin.Text),
			FirstAlbumMax: parseVal(entryAlbumMax.Text),
		}

		filtresActuels.MembersCounts = nil
		for i, chk := range []*widget.Check{chk1, chk2, chk3, chk4} {
			if chk.Checked {
				if i == 3 {
					filtresActuels.MembersCounts = append(filtresActuels.MembersCounts, 4, 5, 6, 7, 8)
				} else {
					filtresActuels.MembersCounts = append(filtresActuels.MembersCounts, i+1)
				}
			}
		}

		filtresActuels.Locations = nil
		if strings.TrimSpace(entryLocations.Text) != "" {
			for _, loc := range strings.Split(entryLocations.Text, ",") {
				l := strings.ToLower(strings.TrimSpace(loc))
				if l != "" {
					filtresActuels.Locations = append(filtresActuels.Locations, l)
				}
			}
		}

		rafraichirGrille()
	})

	btnReset := widget.NewButton("Réinitialiser", func() {
		filtresActuels = FilterOptions{}
		rechercheActuelle = ""
		entryRecherche.SetText("")
		entryCreationMin.SetText("")
		entryCreationMax.SetText("")
		entryAlbumMin.SetText("")
		entryAlbumMax.SetText("")
		entryLocations.SetText("")
		chk1.SetChecked(false)
		chk2.SetChecked(false)
		chk3.SetChecked(false)
		chk4.SetChecked(false)
		AfficherLeMenu()
	})

	colonneFiltres := container.NewVBox(
		widget.NewLabelWithStyle("Filtres", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabel("Date de création :"),
		entryCreationMin,
		entryCreationMax,
		widget.NewSeparator(),
		widget.NewLabel("Premier album (année) :"),
		entryAlbumMin,
		entryAlbumMax,
		widget.NewSeparator(),
		widget.NewLabel("Membres :"),
		chk1, chk2, chk3, chk4,
		widget.NewSeparator(),
		widget.NewLabel("Lieux de concerts :"),
		entryLocations,
		widget.NewSeparator(),
		btnAppliquer,
		btnReset,
	)

	rafraichirGrille()

	split := container.NewHSplit(container.NewVScroll(colonneFiltres), conteneurDefilant)
	split.SetOffset(0.20)

	titrePrincipal := canvas.NewText("Artistes", color.White)
	titrePrincipal.TextSize = 40
	titrePrincipal.TextStyle = fyne.TextStyle{Bold: true}
	titrePrincipal.Alignment = fyne.TextAlignCenter

	contenuCentre := container.NewBorder(
		container.NewVBox(titrePrincipal, barreRecherche),
		nil, nil, nil,
		split,
	)

	maFenetre.SetContent(contenuCentre)
}

func AfficherLesDetails(artiste List_artist) {
	boutonHome := widget.NewButtonWithIcon("Retour", theme.HomeIcon(), func() {
		AfficherLeMenu()
	})

	barreNavigation := container.NewHBox(layout.NewSpacer(), boutonHome)

	var img *canvas.Image
	if uri, err := storage.ParseURI(artiste.Image); err == nil {
		img = canvas.NewImageFromURI(uri)
	} else {
		img = canvas.NewImageFromResource(theme.MediaMusicIcon())
	}
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(350, 350))

	nomTxt := canvas.NewText(artiste.Name, color.Black)
	nomTxt.TextSize = 30
	nomTxt.TextStyle = fyne.TextStyle{Bold: true}

	dateLbl := widget.NewLabelWithStyle(
		"Date de création : "+strconv.Itoa(artiste.CreationDate),
		fyne.TextAlignLeading, fyne.TextStyle{Italic: true},
	)

	albumLbl := widget.NewLabel("Premier Album : " + artiste.FirstAlbum)
	albumLbl.TextStyle = fyne.TextStyle{Bold: true}

	membresTitre := widget.NewLabelWithStyle("Membres :", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Underline: true})
	membresListe := widget.NewLabel("- " + strings.Join(artiste.Members, "\n- "))

	relationsTitre := widget.NewLabelWithStyle("Concerts (Lieux & Dates) :", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Underline: true})

	conteneurRelations := container.NewVBox()

	relations, errRel := GetRelations(artiste.RelationsUrl)
	if errRel != nil {
		conteneurRelations.Add(widget.NewLabel("Impossible de charger les infos concerts."))
	} else {
		for lieu, dates := range relations {
			lieuPropre := strings.Title(strings.ReplaceAll(lieu, "_", " "))
			btn := widget.NewButton(lieuPropre, func(l string) func() {
				return func() { AfficherCarte(l) }
			}(lieuPropre))
			btn.Importance = widget.HighImportance
			lbl := widget.NewLabel(strings.Join(dates, ", "))
			lbl.Wrapping = fyne.TextWrapWord
			conteneurRelations.Add(btn)
			conteneurRelations.Add(lbl)
			conteneurRelations.Add(layout.NewSpacer())
		}
	}

	infosContainer := container.NewVBox(
		nomTxt,
		widget.NewSeparator(),
		dateLbl,
		albumLbl,
		widget.NewSeparator(),
		membresTitre,
		membresListe,
		widget.NewSeparator(),
		relationsTitre,
		conteneurRelations,
	)

	split := container.NewHSplit(img, container.NewVScroll(infosContainer))
	split.SetOffset(0.4)

	pageDetails := container.NewBorder(barreNavigation, nil, nil, nil, split)

	maFenetre.SetContent(pageDetails)

	if desk, ok := maFenetre.Canvas().(desktop.Canvas); ok {
		desk.SetOnKeyDown(func(ev *fyne.KeyEvent) {
			if ev.Name == fyne.KeyEscape {
				AfficherLeMenu()
			}
		})
	}
}

func AfficherCarte(lieu string) {
	lat, lon, nomComplet, err := GeocodeLocation(lieu)
	if err != nil {
		dialog.NewInformation("Localisation", "Impossible de récupérer la géolocalisation.", maFenetre).Show()
		return
	}
	if nomComplet == "" {
		dialog.NewInformation("Localisation", "Aucun résultat trouvé pour ce lieu.", maFenetre).Show()
		return
	}

	const tuileW, tuileH float32 = 256, 256
	const grille = 3
	const mapW = tuileW * grille
	const mapH = tuileH * grille

	zoom := 5

	containerMap := container.NewWithoutLayout()
	containerMap.Resize(fyne.NewSize(mapW, mapH))

	xTuileCenter, yTuileCenter := calculerIndiceTuile(lat, lon, zoom)

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			xTuile := xTuileCenter + dx
			yTuile := yTuileCenter + dy
			mapURL := fmt.Sprintf("https://tile.openstreetmap.org/%d/%d/%d.png", zoom, xTuile, yTuile)
			tuile, err := chargerImageCarteAvecRetry(mapURL)

			if err == nil {
				tuile.Resize(fyne.NewSize(tuileW, tuileH))
				posX := float32(dx+1) * tuileW
				posY := float32(dy+1) * tuileH
				tuile.Move(fyne.NewPos(posX, posY))
				containerMap.Add(tuile)
			}
		}
	}

	xPrecis, yPrecis := calculerPositionPixel(lat, lon, zoom, xTuileCenter, yTuileCenter, tuileW, tuileH, grille)

	marqueur := canvas.NewCircle(color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	marqueur.Resize(fyne.NewSize(14, 14))

	if xPrecis >= 0 && xPrecis <= mapW && yPrecis >= 0 && yPrecis <= mapH {
		marqueur.Move(fyne.NewPos(xPrecis-7, yPrecis-7))
		containerMap.Add(marqueur)
	}

	info := widget.NewLabel(nomComplet)
	info.Wrapping = fyne.TextWrapWord
	info.TextStyle = fyne.TextStyle{Bold: true}

	scrollMap := container.NewVScroll(container.NewHScroll(containerMap))
	scrollMap.SetMinSize(fyne.NewSize(800, 500))

	contenu := container.NewVBox(info, scrollMap)
	dialog.NewCustom("Localisation - OpenStreetMap", "Fermer", contenu, maFenetre).Show()
}

func calculerIndiceTuile(lat, lon float64, zoom int) (int, int) {
	divisor := uint64(1) << uint(zoom)
	x := int((lon + 180.0) / 360.0 * float64(divisor))
	y := int((1 - math.Log(math.Tan(lat*math.Pi/180)+1/math.Cos(lat*math.Pi/180))/math.Pi) / 2 * float64(divisor))
	return x, y
}

func calculerPositionPixel(lat, lon float64, zoom int, xCenter, yCenter int, tuileW, tuileH float32, grille float32) (float32, float32) {
	n := math.Pow(2, float64(zoom))

	xExact := (lon + 180.0) / 360.0 * n
	latRad := lat * math.Pi / 180.0
	yExact := (1.0 - math.Log(math.Tan(latRad)+1.0/math.Cos(latRad))/math.Pi) / 2.0 * n

	xCenterFloat := float64(xCenter)
	yCenterFloat := float64(yCenter)

	xDelta := (xExact - xCenterFloat) * 256.0
	yDelta := (yExact - yCenterFloat) * 256.0

	centerOffsetX := (float32(grille) - 1.0) / 2.0 * tuileW
	centerOffsetY := (float32(grille) - 1.0) / 2.0 * tuileH

	xPixel := centerOffsetX + float32(xDelta)
	yPixel := centerOffsetY + float32(yDelta)

	return xPixel, yPixel
}

func chargerImageCarteAvecRetry(url string) (*canvas.Image, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "image/webp,image/apng,image/svg+xml,image/*,*/*")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("statut %d", resp.StatusCode)
	}

	img := canvas.NewImageFromReader(resp.Body, url)
	if img == nil {
		return nil, fmt.Errorf("image invalide")
	}
	return img, nil
}

type ImageCliquable struct {
	widget.BaseWidget
	Image *canvas.Image
	OnTap func()
}

func NewImageCliquable(img *canvas.Image, onTap func()) *ImageCliquable {
	ic := &ImageCliquable{Image: img, OnTap: onTap}
	ic.ExtendBaseWidget(ic)
	return ic
}

func (i *ImageCliquable) Tapped(_ *fyne.PointEvent) {
	if i.OnTap != nil {
		i.OnTap()
	}
}

func (i *ImageCliquable) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(i.Image)
}

func yearFromDateString(d string) int {
	if len(d) < 4 {
		return 0
	}
	y, err := strconv.Atoi(d[:4])
	if err != nil {
		return 0
	}
	return y
}

func FiltrerArtistes(artistes []List_artist, opts FilterOptions) []List_artist {
	var result []List_artist
	for _, a := range artistes {
		if (opts.CreationMin != 0 && a.CreationDate < opts.CreationMin) ||
			(opts.CreationMax != 0 && a.CreationDate > opts.CreationMax) {
			continue
		}

		if opts.FirstAlbumMin != 0 || opts.FirstAlbumMax != 0 {
			year := yearFromDateString(a.FirstAlbum)
			if year == 0 {
				continue
			}
			if (opts.FirstAlbumMin != 0 && year < opts.FirstAlbumMin) ||
				(opts.FirstAlbumMax != 0 && year > opts.FirstAlbumMax) {
				continue
			}
		}

		if len(opts.MembersCounts) > 0 {
			found := false
			for _, v := range opts.MembersCounts {
				if len(a.Members) == v {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if len(opts.Locations) > 0 {
			rel, err := GetRelations(a.RelationsUrl)
			if err != nil {
				continue
			}
			matchLoc := false
			for loc := range rel {
				locLower := strings.ToLower(strings.ReplaceAll(loc, "_", " "))
				for _, wanted := range opts.Locations {
					if strings.Contains(locLower, wanted) {
						matchLoc = true
						break
					}
				}
				if matchLoc {
					break
				}
			}
			if !matchLoc {
				continue
			}
		}

		result = append(result, a)
	}
	return result
}

func RechercherArtistes(artistes []List_artist, recherche string) []List_artist {
	if recherche == "" {
		return artistes
	}

	rechercheLower := strings.ToLower(recherche)
	var result []List_artist

	for _, a := range artistes {
		if strings.Contains(strings.ToLower(a.Name), rechercheLower) {
			result = append(result, a)
		} else if strings.Contains(strings.ToLower(a.FirstAlbum), rechercheLower) {
			result = append(result, a)
		} else if strings.Contains(strconv.Itoa(a.CreationDate), recherche) {
			result = append(result, a)
		} else {
			found := false
			for _, membre := range a.Members {
				if strings.Contains(strings.ToLower(membre), rechercheLower) {
					result = append(result, a)
					found = true
					break
				}
			}
			if found {
				continue
			}

			rel, err := GetRelations(a.RelationsUrl)
			if err == nil {
				for loc := range rel {
					locLower := strings.ToLower(strings.ReplaceAll(loc, "_", " "))
					if strings.Contains(locLower, rechercheLower) {
						result = append(result, a)
						break
					}
				}
			}
		}
	}

	return result
}
