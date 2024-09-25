package models

type Patient struct {
	ID                     int     `json:"id"`
	DateOfBaselineVisit    string  `json:"date_of_baseline_visit"`
	AgeVisitBaseline       int     `json:"age_visit_baseline"`
	HypertensionBaseline   bool    `json:"hypertension_baseline"`
	DiabetesBaseline       bool    `json:"diabetes_baseline"`
	SmokingStatusBaseline  string  `json:"smoking_status_baseline"`
	CVDBaseline            bool    `json:"cvd_baseline"`
	CADBaseline            bool    `json:"cad_baseline"`
	MIBaseline             bool    `json:"mi_baseline"`
	CADRevascularization   bool    `json:"cad_revascularization_baseline"`
	StrokeBaseline         bool    `json:"stroke_baseline"`
	StrokePremature        bool    `json:"stroke_premature_baseline"`
	LiverSteatosisBaseline bool    `json:"liver_steatosis_baseline"`
	Xanthelasma            bool    `json:"xanthelasma"`
	WeightBaseline         float64 `json:"weight_baseline"`
	HeightBaseline         float64 `json:"height_baseline"`
	ThyroidDisease         bool    `json:"thyroid_disease"`
	MenarcheBaseline       bool    `json:"menarche_reached_baseline"`
	AgeMenarche            int     `json:"age_at_menarche_baseline"`
	MenopauseBaseline      bool    `json:"menopause_reached_baseline"`
}
