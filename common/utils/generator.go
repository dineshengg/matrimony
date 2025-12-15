package utils

import (
	"fmt"
	"sync"
	"unicode"

	log "github.com/sirupsen/logrus"
)

type SafeCounter struct {
	//global counter to know the total users enrolled in our site
	counter      uint64
	mu           sync.Mutex
	mid          map[string]uint64
	langmap      map[string]byte
	religionmap  map[string]byte
	communitymap map[string]byte
}

var matCounter SafeCounter

func init() {
	log.Debug("Initializing matrimony ids creation")
}
func IsMatrimonyIDInit() bool {
	matCounter.mu.Lock()
	defer matCounter.mu.Unlock()
	if matCounter.mid == nil {
		matCounter.mid = make(map[string]uint64)
	}
	if matCounter.langmap == nil {
		matCounter.langmap = make(map[string]byte)
	}
	if matCounter.religionmap == nil {
		matCounter.religionmap = make(map[string]byte)
	}
	if matCounter.communitymap == nil {
		matCounter.communitymap = make(map[string]byte)
	}

	for _, lang := range []string{"tamil", "telugu", "malayalam", "kannada", "hindi", "english", "others"} {
		matCounter.langmap[lang] = byte(unicode.ToUpper(rune(lang[0])))
	}
	for _, community := range []string{"hindu", "muslim", "christian", "sikh", "jain", "others"} {
		matCounter.communitymap[community] = byte(unicode.ToUpper(rune(community[0])))
	}
	for _, religion := range []string{"hindu", "muslim", "christian", "sikh", "others"} {
		matCounter.religionmap[religion] = byte(unicode.ToUpper(rune(religion[0])))
	}

	for _, val := range matCounter.langmap {
		for _, rel := range matCounter.religionmap {
			for _, community := range matCounter.communitymap {
				Matid := fmt.Sprintf("KAN%c%c%c", val, rel, community)
				matCounter.mid[Matid] = 0
			}
		}
	}

	matCounter.mid["global"] = matCounter.counter

	db := GetDB()
	if db == nil {
		log.Error("Not able to get database connection")
		return false
	}

	if err := db.Raw("select counter from globalcounters order by counter desc limit 1").Scan(&matCounter.counter); err != nil {
		log.Errorf("Error getting global counter - %v", err)
		return false
	}
	return true

}
func StoreMatrimonyIDCounters() error {
	matCounter.mu.Lock()
	defer matCounter.mu.Unlock()

	for matrimonyID, counter := range matCounter.mid {
		log.Debugf("Storing Matrimony ID: %s with count: %d", matrimonyID, counter)
		if err := GetDB().Exec("INSERT INTO globalcounters (category, counter) values (?, ?) ON CONFLICT (category) DO UPDATE SET category = EXCLUDED.category, counter = EXCULDED.counter", matrimonyID, counter).Error; err != nil {
			log.Errorf("Error storing Matrimony ID %s: %v", matrimonyID, err)
			return fmt.Errorf("error storing Matrimony ID %s: %v", matrimonyID, err)
		}
	}

	return nil
}

func MatrimonyID(lang, religion, community string) string {
	matCounter.mu.Lock()
	defer matCounter.mu.Unlock()

	//currently we are using mat id which is common for all users
	var MatID string
	if len(lang) == 0 || len(religion) == 0 || len(community) == 0 {
		log.Info("generating common matrimony id for all users")

		MatID = fmt.Sprintf("KAN%06d", matCounter.counter)
		return MatID
	}

	var l, r, c byte
	var ok bool
	//MATID - KAN<lang 1 byte><religon 1 byte><community 1 byte><counter 6 bytes>
	if l, ok = matCounter.langmap[lang]; !ok {
		log.Infof("Language not found %s", lang)
		l = 'u'
	}
	if r, ok = matCounter.religionmap[religion]; !ok {
		log.Infof("religion not found %s", religion)
		r = 'u'
	}
	if c, ok = matCounter.communitymap[community]; !ok {
		log.Infof("community not found %s", community)
		c = 'u'
	}

	MatID_Pre := fmt.Sprintf("KAN%c%c%c", l, r, c)
	if _, ok := matCounter.mid[MatID_Pre]; !ok {
		log.Errorf("Matrimony id prefix is not in category using default counter %s", MatID_Pre)
		MatID = fmt.Sprintf("KAN%06d", matCounter.counter)
	} else {
		catCounter := matCounter.mid[MatID_Pre]
		MatID = fmt.Sprintf("%s%06d", MatID_Pre, catCounter)
		matCounter.mid[MatID_Pre] = catCounter + 1
	}
	matCounter.counter++

	log.Debugf("Matrimony ID generated: %s", MatID)
	return MatID

}
