# driver-bbc-iplayer

A databox driver to BBC iPlayer recommendations


## status

work in progress

## Data

This driver will update the users recommendations every hour and store them im a KVJSON store datasource: IplayerRecommend key: all.

The data has the following format:

https://ibl.api.bbci.co.uk/ibl/v1/schema/ibl.json under user_recommendations

example:

```JSON
{
	"version": "1.0",
	"schema": "/ibl/v1/schema/ibl.json",
	"user_recommendations": {
		"rec_source": "cw",
		"rec_set": "305.16619.1541670734337",
		"rec_feed": "ipew-boxc",
		"elements": [{
			"type": "user_recommendation",
			"algorithm": "cw-fallback-ipew-boxc",
			"episode": {
				"id": "p06crp3c",
				"live": false,
				"type": "episode_large",
				"title": "Bodyguard",
				"images": {
					"type": "image",
					"standard": "https://ichef.bbci.co.uk/images/ic/{recipe}/p06j4zb7.jpg"
				},
				"labels": {
					"category": "Drama"
				},
				"signed": false,
				"status": "available",
				"tleo_id": "p06crngy",
				"guidance": true,
				"subtitle": "Series 1: Episode 1",
				"synopses": {
					"large": "After distinguishing himself by courageously neutralising a terrorist threat, troubled war veteran Police Sergeant David Budd of the Metropolitan Police's Royalty and Specialist Protection Branch (RaSP) is assigned as a principal protection officer to the home secretary, the Rt Hon Julia Montague MP.\n\nJulia is a controversial politician intent on pushing a new counterterrorism bill through Parliament which would give the security service enhanced surveillance powers. Her political ambitions make Julia a high-profile target. Dedicated to his job but concealing deep resentment of politicians following his traumas in Helmand, David's divided loyalties might make him Julia's greatest threat.",
					"small": "A troubled war veteran is assigned to protect a controversial politician.",
					"medium": "Drama series. A troubled war veteran, PS David Budd, is assigned to protect a controversial politician who may be the target of a terror plot.",
					"editorial": "From the writer of Line of Duty. A war veteran is torn between his duty and his beliefs."
				},
				"versions": [{
					"hd": true,
					"id": "b0bhr5wr",
					"uhd": false,
					"kind": "original",
					"type": "version_large",
					"events": [{
						"name": "started",
						"offset": 30,
						"system": "uas"
					}, {
						"name": "ended",
						"offset": 3361,
						"system": "uas"
					}, {
						"name": "iplxp-ep-started",
						"offset": 30,
						"system": "optimizely"
					}, {
						"name": "iplxp-ep-watched",
						"offset": 3106,
						"system": "optimizely"
					}, {
						"name": "iplxp-ep-started",
						"offset": 30,
						"system": "dax"
					}, {
						"name": "iplxp-ep-watched",
						"offset": 3106,
						"system": "dax"
					}],
					"download": true,
					"duration": {
						"text": "58 mins",
						"value": "PT57M31.080S"
					},
					"guidance": {
						"id": "D2L1",
						"text": {
							"large": "Contains some strong language and scenes which some viewers may find upsetting.",
							"small": "Has guidance",
							"medium": "Contains some strong language and upsetting scenes."
						}
					},
					"availability": {
						"end": "2019-03-23T21:00:00Z",
						"start": "2018-08-26T21:00:00Z",
						"remaining": {
							"text": "Available for 4 months"
						}
					},
					"credits_start": 3408,
					"first_broadcast": "9pm 26 Aug 2018"
				}, {
					"hd": false,
					"id": "p06jjjyx",
					"uhd": false,
					"kind": "audio-described",
					"type": "version_large",
					"events": [{
						"name": "started",
						"offset": 30,
						"system": "uas"
					}, {
						"name": "ended",
						"offset": 3361,
						"system": "uas"
					}, {
						"name": "iplxp-ep-started",
						"offset": 30,
						"system": "optimizely"
					}, {
						"name": "iplxp-ep-watched",
						"offset": 3106,
						"system": "optimizely"
					}, {
						"name": "iplxp-ep-started",
						"offset": 30,
						"system": "dax"
					}, {
						"name": "iplxp-ep-watched",
						"offset": 3106,
						"system": "dax"
					}],
					"download": false,
					"duration": {
						"text": "58 mins",
						"value": "PT57M30.600S"
					},
					"guidance": {
						"id": "D2L1",
						"text": {
							"large": "Contains some strong language and scenes which some viewers may find upsetting.",
							"small": "Has guidance",
							"medium": "Contains some strong language and upsetting scenes."
						}
					},
					"availability": {
						"end": "2019-03-23T21:00:00Z",
						"start": "2018-08-26T21:39:07Z",
						"remaining": {
							"text": "Available for 4 months"
						}
					},
					"first_broadcast": "26 Aug 2018"
				}],
				"parent_id": "p06crnq6",
				"tleo_type": "brand",
				"categories": ["audio-described", "drama-and-soaps"],
				"has_credits": true,
				"master_brand": {
					"id": "bbc_one",
					"titles": {
						"large": "BBC One",
						"small": "BBC One",
						"medium": "BBC One"
					},
					"ident_id": "p06g87jc",
					"attribution": "bbc_one"
				},
				"release_date": "26 Aug 2018",
				"related_links": [{
					"id": "p06jnjc7",
					"url": "http://www.bbc.co.uk/blogs/writersroom/entries/f5c09a37-a04d-4238-b4b3-0ba8f23a9a5d",
					"kind": "priority_content",
					"type": "link",
					"title": "BBC Writersroom: Writer Jed Mercurio introduces Bodyguard"
				}],
				"original_title": "Episode 1",
				"programme_type": "narrative",
				"audio_described": true,
				"parent_position": 1,
				"requires_sign_in": true,
				"release_date_time": "2018-08-26T00:00:00.000Z",
				"editorial_subtitle": "Box Set. He has to keep her safe. Is he her biggest threat?",
				"lexical_sort_letter": "B",
				"requires_tv_licence": true
			}
		}, {
			"type": "user_recommendation",
			"algorithm": "cw-fallback-ipew-boxc",
			"episode": {
				"id": "p06kbg8t",
				"live": false,
				"type": "episode_large",
				"tests": [{
					"id": "killingeve_S1E1",
					"variants": [{
						"id": "variant1",
						"data": {
							"images": {
								"type": "image",
								"standard": "https://ichef.bbci.co.uk/images/ic/{recipe}/p06kymk7.jpg"
							}
						}
					}, {
						"id": "defaultvariant",
						"data": {
							"images": {
								"type": "image",
								"standard": "https://ichef.bbci.co.uk/images/ic/{recipe}/p06kym8w.jpg"
							}
						}
					}]
				}],
				"title": "Killing Eve",
				"images": {
					"type": "image",
					"standard": "https://ichef.bbci.co.uk/images/ic/{recipe}/p06kym8w.jpg",
					"programme_logo": "https://ichef.bbci.co.uk/images/ic/{recipe}/p06ng7vj.jpg"
				},
				"labels": {
					"category": "Drama"
				},
				"signed": false,
				"status": "available",
				"tleo_id": "p06jy6bc",
				"guidance": true,
				"subtitle": "Series 1: 1. Nice Face",
				"synopses": {
					"large": "Eve is a security officer at MI5. She likes her job, her boss, her friends and her husband, but she is bored. When a Russian diplomat is assassinated in Vienna and Eve is given the job of looking after the witness, a casual bet about the identity of the killer gets out of hand, and Eve finds herself drawn into a cat-and-mouse game all across the continent, looking for a deadly, elusive and fascinating suspect.",
					"small": "When a politician is murdered, an MI5 security officer must protect the only witness.",
					"medium": "Thriller in which a security operative hunts for an assassin. When a Russian politician is murdered, MI5 security officer Eve must protect the only witness.",
					"preview": "Two women go head to head in an epic game of cat and mouse.",
					"editorial": "Eve craves excitement. Villanelle craves destruction. They both need each other."
				},
				"versions": [{
					"hd": true,
					"id": "p06kbh1r",
					"uhd": false,
					"kind": "original",
					"type": "version_large",
					"events": [{
						"name": "started",
						"offset": 30,
						"system": "uas"
					}, {
						"name": "ended",
						"offset": 2457,
						"system": "uas"
					}, {
						"name": "iplxp-ep-started",
						"offset": 30,
						"system": "optimizely"
					}, {
						"name": "iplxp-ep-watched",
						"offset": 2292,
						"system": "optimizely"
					}, {
						"name": "iplxp-ep-started",
						"offset": 30,
						"system": "dax"
					}, {
						"name": "iplxp-ep-watched",
						"offset": 2292,
						"system": "dax"
					}],
					"download": true,
					"duration": {
						"text": "42 mins",
						"value": "PT42M26.960S"
					},
					"guidance": {
						"id": "D1V1",
						"text": {
							"large": "Contains some violent scenes and some scenes which some viewers may find upsetting.",
							"small": "Has guidance",
							"medium": "Contains some violence and some upsetting scenes."
						}
					},
					"availability": {
						"end": "2022-08-31T22:59:00Z",
						"start": "2018-09-15T20:15:00Z",
						"remaining": {
							"text": "Available for over a year"
						}
					},
					"credits_start": 2519,
					"first_broadcast": "9:15pm 15 Sep 2018"
				}, {
					"hd": false,
					"id": "p06l9fwc",
					"uhd": false,
					"kind": "audio-described",
					"type": "version_large",
					"events": [{
						"name": "started",
						"offset": 30,
						"system": "uas"
					}, {
						"name": "ended",
						"offset": 2458,
						"system": "uas"
					}, {
						"name": "iplxp-ep-started",
						"offset": 30,
						"system": "optimizely"
					}, {
						"name": "iplxp-ep-watched",
						"offset": 2293,
						"system": "optimizely"
					}, {
						"name": "iplxp-ep-started",
						"offset": 30,
						"system": "dax"
					}, {
						"name": "iplxp-ep-watched",
						"offset": 2293,
						"system": "dax"
					}],
					"download": false,
					"duration": {
						"text": "42 mins",
						"value": "PT42M27.968S"
					},
					"guidance": {
						"id": "D1V1",
						"text": {
							"large": "Contains some violent scenes and some scenes which some viewers may find upsetting.",
							"small": "Has guidance",
							"medium": "Contains some violence and some upsetting scenes."
						}
					},
					"availability": {
						"end": "2022-08-31T22:59:00Z",
						"start": "2018-09-15T21:44:38Z",
						"remaining": {
							"text": "Available for over a year"
						}
					},
					"first_broadcast": "15 Sep 2018"
				}],
				"parent_id": "p06jy6gl",
				"tleo_type": "brand",
				"categories": ["audio-described", "drama-and-soaps"],
				"preview_id": "p06p3ftx",
				"has_credits": true,
				"master_brand": {
					"id": "bbc_three",
					"titles": {
						"large": "BBC Three",
						"small": "BBC Three",
						"medium": "BBC Three"
					},
					"ident_id": "p03bhqhy",
					"attribution": "bbc_three"
				},
				"release_date": "15 Sep 2018",
				"related_links": [],
				"original_title": "Nice Face",
				"programme_type": "narrative",
				"audio_described": true,
				"parent_position": 1,
				"requires_sign_in": true,
				"release_date_time": "2018-09-15T00:00:00.000Z",
				"editorial_subtitle": "Box Set. A deadly, obsessive game of cat and mouse",
				"lexical_sort_letter": "K",
				"requires_tv_licence": true
			}
		}]
	}
}
```

## Development of databox was supported by the following funding
```
EP/N028260/1, Databox: Privacy-Aware Infrastructure for Managing Personal Data
EP/N028260/2, Databox: Privacy-Aware Infrastructure for Managing Personal Data
EP/N014243/1, Future Everyday Interaction with the Autonomous Internet of Things
EP/M001636/1, Privacy-by-Design: Building Accountability into the Internet of Things EP/M02315X/1, From Human Data to Personal Experience
```
