package huddles

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

type HuddleProfileSuite struct {
	Huddle      *models.Group
	HuddleBSONM bson.M
	HuddleJSON  []byte
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&HuddleProfileSuite{})

func (h *HuddleProfileSuite) SetUpTest(c *C) {
	h.Huddle = newExampleHuddle()
	h.HuddleBSONM = newExampleHuddleBSONM()

	data, err := ioutil.ReadFile("../fixtures/huddle.json")
	util.CheckErr(err)
	h.HuddleJSON = data
}

func (h *HuddleProfileSuite) TearDownTest(c *C) {
	h.Huddle = nil
	h.HuddleBSONM = bson.M{}
	h.HuddleJSON = make([]byte, 0)
}

func (h *HuddleProfileSuite) TestMarshalJSON(c *C) {
	// First unmarshal the expected JSON into a map for easier comparison
	var expected map[string]interface{}
	err := json.Unmarshal(h.HuddleJSON, &expected)
	util.CheckErr(err)

	// Now marshal the huddle struct into JSON data
	data, err := json.Marshal(h.Huddle)
	util.CheckErr(err)

	// In order to compare it, unmarshal the new JSON data to a map
	var obtained map[string]interface{}
	err = json.Unmarshal(data, &obtained)
	util.CheckErr(err)

	c.Assert(obtained, DeepEquals, expected)
}

func (h *HuddleProfileSuite) TestUnmarshalJSON(c *C) {
	obtained := &models.Group{}
	err := json.Unmarshal(h.HuddleJSON, &obtained)
	util.CheckErr(err)

	assertDeepEqualHuddles(c, obtained, h.Huddle)
}

func (h *HuddleProfileSuite) TestRoundTripJSON(c *C) {
	// First we need to marshal it, just so we can get our bytes to unmarshal
	data, err := json.Marshal(h.Huddle)
	util.CheckErr(err)

	// Now we'll try unmarshalling.  If everything is working it should survive the round trip.
	var obtained models.Group
	err = json.Unmarshal(data, &obtained)
	util.CheckErr(err)

	assertDeepEqualHuddles(c, &obtained, h.Huddle)
}

func (h *HuddleProfileSuite) TestMarshalBSON(c *C) {
	// Now marshal the huddle struct into BSON data
	data, err := bson.Marshal(h.Huddle)
	util.CheckErr(err)

	// In order to compare it, unmarshal the new BSON data to a map
	var obtained bson.M
	err = bson.Unmarshal(data, &obtained)
	util.CheckErr(err)

	// Since times don't work in DeepEquals (due to timezoney shenanigans in Go), first check the times directly.
	// After confirming they represent the same moment, set the expected to the obtained so we pass DeepEquals
	c.Assert(obtained["extension"].([]interface{})[0].(bson.M)["activeDateTime"].(bson.M)["time"].(time.Time).Unix(), Equals,
		h.HuddleBSONM["extension"].([]interface{})[0].(bson.M)["activeDateTime"].(bson.M)["time"].(time.Time).Unix())
	h.HuddleBSONM["extension"].([]interface{})[0].(bson.M)["activeDateTime"].(bson.M)["time"] =
		obtained["extension"].([]interface{})[0].(bson.M)["activeDateTime"].(bson.M)["time"]
	for i := 0; i < 3; i++ {
		c.Assert(obtained["member"].([]interface{})[i].(bson.M)["extension"].([]interface{})[1].(bson.M)["reviewed"].(bson.M)["time"].(time.Time).Unix(), Equals,
			h.HuddleBSONM["member"].([]interface{})[i].(bson.M)["extension"].([]interface{})[1].(bson.M)["reviewed"].(bson.M)["time"].(time.Time).Unix())
		h.HuddleBSONM["member"].([]interface{})[i].(bson.M)["extension"].([]interface{})[1].(bson.M)["reviewed"].(bson.M)["time"] =
			obtained["member"].([]interface{})[i].(bson.M)["extension"].([]interface{})[1].(bson.M)["reviewed"].(bson.M)["time"]
	}

	c.Assert(obtained, DeepEquals, h.HuddleBSONM)
}

func (h *HuddleProfileSuite) TestUnmarshalBSON(c *C) {
	// First we need to marshal the expected BSON into BSON data
	data, err := bson.Marshal(h.HuddleBSONM)
	util.CheckErr(err)

	// Then unmarshal the BSON data into the Huddle struct
	obtained := &models.Group{}
	err = bson.Unmarshal(data, &obtained)
	util.CheckErr(err)

	assertDeepEqualHuddles(c, obtained, h.Huddle)
}

func (h *HuddleProfileSuite) TestRoundTripBSON(c *C) {
	// First we need to marshal it, just so we can get our bytes to unmarshal
	data, err := bson.Marshal(h.Huddle)
	util.CheckErr(err)

	// Now we'll try unmarshalling.  If everything is working it should survive the round trip.
	var obtained models.Group
	err = bson.Unmarshal(data, &obtained)
	util.CheckErr(err)

	assertDeepEqualHuddles(c, &obtained, h.Huddle)
}

func newExampleHuddle() *models.Group {
	tru := true
	return &models.Group{
		DomainResource: models.DomainResource{
			Resource: models.Resource{
				ResourceType: "Group",
				Meta: &models.Meta{
					Profile: []string{"http://interventionengine.org/fhir/profile/huddle"},
				},
			},
			Extension: []models.Extension{
				{
					Url:           "http://interventionengine.org/fhir/extension/group/activeDateTime",
					ValueDateTime: &models.FHIRDateTime{Time: time.Date(2016, time.February, 2, 9, 0, 0, 0, time.UTC), Precision: models.Precision(models.Timestamp)},
				},
				{
					Url: "http://interventionengine.org/fhir/extension/group/leader",
					ValueReference: &models.Reference{
						Reference:    "Practitioner/9999999999999999999",
						ReferencedID: "9999999999999999999",
						Type:         "Practitioner",
						External:     new(bool),
					},
				},
			},
		},
		Type:   "person",
		Actual: &tru,
		Code: &models.CodeableConcept{
			Coding: []models.Coding{
				{System: "http://interventionengine.org/fhir/cs/huddle", Code: "HUDDLE"},
			},
			Text: "Huddle",
		},
		Name: "Dr. Smith's Huddle for February 2, 2016",
		Member: []models.GroupMemberComponent{
			{
				BackboneElement: models.BackboneElement{
					Element: models.Element{
						Extension: []models.Extension{
							{
								Url: "http://interventionengine.org/fhir/extension/group/member/reason",
								ValueCodeableConcept: &models.CodeableConcept{
									Coding: []models.Coding{
										{System: "http://interventionengine.org/fhir/cs/huddle-member-reason", Code: "RECENT_ADMISSION"},
									},
									Text: "Recent Inpatient Admission",
								},
							},
							{
								Url:           "http://interventionengine.org/fhir/extension/group/member/reviewed",
								ValueDateTime: &models.FHIRDateTime{Time: time.Date(2016, time.February, 2, 9, 8, 15, 0, time.UTC), Precision: models.Precision(models.Timestamp)},
							},
						},
					},
				},
				Entity: &models.Reference{
					Reference:    "Patient/1111111111111111111",
					ReferencedID: "1111111111111111111",
					Type:         "Patient",
					External:     new(bool),
				},
			},
			{
				BackboneElement: models.BackboneElement{
					Element: models.Element{
						Extension: []models.Extension{
							{
								Url: "http://interventionengine.org/fhir/extension/group/member/reason",
								ValueCodeableConcept: &models.CodeableConcept{
									Coding: []models.Coding{
										{System: "http://interventionengine.org/fhir/cs/huddle-member-reason", Code: "RISK_SCORE"},
									},
									Text: "Risk Score Warrants Discussion",
								},
							},
							{
								Url:           "http://interventionengine.org/fhir/extension/group/member/reviewed",
								ValueDateTime: &models.FHIRDateTime{Time: time.Date(2016, time.February, 2, 9, 15, 46, 0, time.UTC), Precision: models.Precision(models.Timestamp)},
							},
						},
					},
				},
				Entity: &models.Reference{
					Reference:    "Patient/2222222222222222222",
					ReferencedID: "2222222222222222222",
					Type:         "Patient",
					External:     new(bool),
				},
			},
			{
				BackboneElement: models.BackboneElement{
					Element: models.Element{
						Extension: []models.Extension{
							{
								Url: "http://interventionengine.org/fhir/extension/group/member/reason",
								ValueCodeableConcept: &models.CodeableConcept{
									Coding: []models.Coding{
										{System: "http://interventionengine.org/fhir/cs/huddle-member-reason", Code: "MANUAL_ADDITION"},
									},
									Text: "Manually Added",
								},
							},
							{
								Url:           "http://interventionengine.org/fhir/extension/group/member/reviewed",
								ValueDateTime: &models.FHIRDateTime{Time: time.Date(2016, time.February, 2, 9, 32, 15, 0, time.UTC), Precision: models.Precision(models.Timestamp)},
							},
						},
					},
				},
				Entity: &models.Reference{
					Reference:    "Patient/3333333333333333333",
					ReferencedID: "3333333333333333333",
					Type:         "Patient",
					External:     new(bool),
				},
			},
			{
				BackboneElement: models.BackboneElement{
					Element: models.Element{
						Extension: []models.Extension{
							{
								Url: "http://interventionengine.org/fhir/extension/group/member/reason",
								ValueCodeableConcept: &models.CodeableConcept{
									Coding: []models.Coding{
										{System: "http://interventionengine.org/fhir/cs/huddle-member-reason", Code: "RECENT_ED_VISIT"},
									},
									Text: "Recent Emergency Department Visit",
								},
							},
						},
					},
				},
				Entity: &models.Reference{
					Reference:    "Patient/4444444444444444444",
					ReferencedID: "4444444444444444444",
					Type:         "Patient",
					External:     new(bool),
				},
			},
			{
				BackboneElement: models.BackboneElement{
					Element: models.Element{
						Extension: []models.Extension{
							{
								Url: "http://interventionengine.org/fhir/extension/group/member/reason",
								ValueCodeableConcept: &models.CodeableConcept{
									Coding: []models.Coding{
										{System: "http://interventionengine.org/fhir/cs/huddle-member-reason", Code: "RECENT_READMISSION"},
									},
									Text: "Recent Inpatient Readmission",
								},
							},
						},
					},
				},
				Entity: &models.Reference{
					Reference:    "Patient/5555555555555555555",
					ReferencedID: "5555555555555555555",
					Type:         "Patient",
					External:     new(bool),
				},
			},
		},
	}
}

func newExampleHuddleBSONM() bson.M {
	return bson.M{
		"resourceType": "Group",
		"meta": bson.M{
			"profile": []interface{}{"http://interventionengine.org/fhir/profile/huddle"},
		},
		"extension": []interface{}{
			bson.M{
				"@context": bson.M{
					"activeDateTime": bson.M{
						"@id":   "http://interventionengine.org/fhir/extension/group/activeDateTime",
						"@type": "dateTime",
					},
				},
				"activeDateTime": bson.M{
					"time":      time.Date(2016, time.February, 2, 9, 0, 0, 0, time.UTC),
					"precision": "timestamp",
				},
			},
			bson.M{
				"@context": bson.M{
					"leader": bson.M{
						"@id":   "http://interventionengine.org/fhir/extension/group/leader",
						"@type": "Reference",
					},
				},
				"leader": bson.M{
					"reference":   "Practitioner/9999999999999999999",
					"referenceid": "9999999999999999999",
					"type":        "Practitioner",
					"external":    false,
				},
			},
		},
		"type":   "person",
		"actual": true,
		"code": bson.M{
			"coding": []interface{}{
				bson.M{
					"system": "http://interventionengine.org/fhir/cs/huddle",
					"code":   "HUDDLE",
				},
			},
			"text": "Huddle",
		},
		"name": "Dr. Smith's Huddle for February 2, 2016",
		"member": []interface{}{
			bson.M{
				"extension": []interface{}{
					bson.M{
						"@context": bson.M{
							"reason": bson.M{
								"@id":   "http://interventionengine.org/fhir/extension/group/member/reason",
								"@type": "CodeableConcept",
							},
						},
						"reason": bson.M{
							"coding": []interface{}{
								bson.M{
									"system": "http://interventionengine.org/fhir/cs/huddle-member-reason",
									"code":   "RECENT_ADMISSION",
								},
							},
							"text": "Recent Inpatient Admission",
						},
					},
					bson.M{
						"@context": bson.M{
							"reviewed": bson.M{
								"@id":   "http://interventionengine.org/fhir/extension/group/member/reviewed",
								"@type": "dateTime",
							},
						},
						"reviewed": bson.M{
							"time":      time.Date(2016, time.February, 2, 9, 8, 15, 0, time.UTC),
							"precision": "timestamp",
						},
					},
				},
				"entity": bson.M{
					"reference":   "Patient/1111111111111111111",
					"referenceid": "1111111111111111111",
					"type":        "Patient",
					"external":    false,
				},
			},
			bson.M{
				"extension": []interface{}{
					bson.M{
						"@context": bson.M{
							"reason": bson.M{
								"@id":   "http://interventionengine.org/fhir/extension/group/member/reason",
								"@type": "CodeableConcept",
							},
						},
						"reason": bson.M{
							"coding": []interface{}{
								bson.M{
									"system": "http://interventionengine.org/fhir/cs/huddle-member-reason",
									"code":   "RISK_SCORE",
								},
							},
							"text": "Risk Score Warrants Discussion",
						},
					},
					bson.M{
						"@context": bson.M{
							"reviewed": bson.M{
								"@id":   "http://interventionengine.org/fhir/extension/group/member/reviewed",
								"@type": "dateTime",
							},
						},
						"reviewed": bson.M{
							"time":      time.Date(2016, time.February, 2, 9, 15, 46, 0, time.UTC),
							"precision": "timestamp",
						},
					},
				},
				"entity": bson.M{
					"reference":   "Patient/2222222222222222222",
					"referenceid": "2222222222222222222",
					"type":        "Patient",
					"external":    false,
				},
			},
			bson.M{
				"extension": []interface{}{
					bson.M{
						"@context": bson.M{
							"reason": bson.M{
								"@id":   "http://interventionengine.org/fhir/extension/group/member/reason",
								"@type": "CodeableConcept",
							},
						},
						"reason": bson.M{
							"coding": []interface{}{
								bson.M{
									"system": "http://interventionengine.org/fhir/cs/huddle-member-reason",
									"code":   "MANUAL_ADDITION",
								},
							},
							"text": "Manually Added",
						},
					},
					bson.M{
						"@context": bson.M{
							"reviewed": bson.M{
								"@id":   "http://interventionengine.org/fhir/extension/group/member/reviewed",
								"@type": "dateTime",
							},
						},
						"reviewed": bson.M{
							"time":      time.Date(2016, time.February, 2, 9, 32, 15, 0, time.UTC),
							"precision": "timestamp",
						},
					},
				},
				"entity": bson.M{
					"reference":   "Patient/3333333333333333333",
					"referenceid": "3333333333333333333",
					"type":        "Patient",
					"external":    false,
				},
			},
			bson.M{
				"extension": []interface{}{
					bson.M{
						"@context": bson.M{
							"reason": bson.M{
								"@id":   "http://interventionengine.org/fhir/extension/group/member/reason",
								"@type": "CodeableConcept",
							},
						},
						"reason": bson.M{
							"coding": []interface{}{
								bson.M{
									"system": "http://interventionengine.org/fhir/cs/huddle-member-reason",
									"code":   "RECENT_ED_VISIT",
								},
							},
							"text": "Recent Emergency Department Visit",
						},
					},
				},
				"entity": bson.M{
					"reference":   "Patient/4444444444444444444",
					"referenceid": "4444444444444444444",
					"type":        "Patient",
					"external":    false,
				},
			},
			bson.M{
				"extension": []interface{}{
					bson.M{
						"@context": bson.M{
							"reason": bson.M{
								"@id":   "http://interventionengine.org/fhir/extension/group/member/reason",
								"@type": "CodeableConcept",
							},
						},
						"reason": bson.M{
							"coding": []interface{}{
								bson.M{
									"system": "http://interventionengine.org/fhir/cs/huddle-member-reason",
									"code":   "RECENT_READMISSION",
								},
							},
							"text": "Recent Inpatient Readmission",
						},
					},
				},
				"entity": bson.M{
					"reference":   "Patient/5555555555555555555",
					"referenceid": "5555555555555555555",
					"type":        "Patient",
					"external":    false,
				},
			},
		},
	}
}

func assertDeepEqualHuddles(c *C, obtained *models.Group, expected *models.Group) {
	// Since times don't work in DeepEquals (due to timezoney shenanigans in Go), first check the times directly.
	// After confirming they represent the same moment, set the expected to the obtained so we pass DeepEquals
	assertAndFixDeepEqualExtensions(c, obtained.Extension, expected.Extension)
	for i := 0; i < len(obtained.Member) && i < len(expected.Member); i++ {
		assertAndFixDeepEqualExtensions(c, obtained.Member[i].Extension, expected.Member[i].Extension)
	}
	c.Assert(obtained, DeepEquals, expected)
}

func assertAndFixDeepEqualExtensions(c *C, obtained []models.Extension, expected []models.Extension) {
	// Since times don't work in DeepEquals (due to timezoney shenanigans in Go), first check the times directly.
	// After confirming they represent the same moment, set the expected to the obtained so we pass DeepEquals
	for i := 0; i < len(obtained) && i < len(expected); i++ {
		if obtained[i].ValueDateTime != nil && expected[i].ValueDateTime != nil {
			c.Assert(obtained[i].ValueDateTime.Time.Unix(), Equals, expected[i].ValueDateTime.Time.Unix())
			expected[i].ValueDateTime.Time = obtained[i].ValueDateTime.Time
		}
	}
	c.Assert(obtained, DeepEquals, expected)
}
