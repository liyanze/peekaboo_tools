package sensitive

import (
	"cmp"
	"encoding/json"
	ahocorasick "github.com/anknown/ahocorasick"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"
)

type (
	ParseSensitiveWords struct {
	}
)

func newParseSensitiveWords() ParseSensitiveWords {
	return ParseSensitiveWords{}
}

type (
	JsonSensitiveKeyWordsPayload struct {
		Sex      []string `json:"sex"`
		Suicide  []string `json:"suicide"`
		Politics []string `json:"politics"`
		Underage []string `json:"underage"`
	}
)

func buildAhoCorasickAutomaton(sensitiveWords [][]rune) *ahocorasick.Machine {
	m := &ahocorasick.Machine{}
	m.Build(sensitiveWords)
	return m

}
func findSensitiveWords(text string, automaton *ahocorasick.Machine) []string {
	now := time.Now()
	defer func() {
		slog.Info("", time.Now().UnixMilli()-now.UnixMilli())
	}()

	foundWords := make(map[string]bool)
	indices := automaton.MultiPatternSearch([]rune(text), false)

	for _, index := range indices {
		foundWords[string(index.Word)] = true
	}
	result := make([]string, 0, len(foundWords))
	for word := range foundWords {
		result = append(result, word)
	}
	return result
}

func (p ParseSensitiveWords) do() {
	pornWords := p.doDetail(pornKeywordsText)
	var tmpPornWords [][]rune
	for _, pw := range pornWords {
		tmpPornWords = append(tmpPornWords, []rune(pw))
	}
	automaton := buildAhoCorasickAutomaton(tmpPornWords)
	words := findSensitiveWords("sex is a bad words, fucking sex, penis", automaton)
	slog.Info("", words)
	suicideWords := p.doDetail(suicideKeywordsText)
	politicsWords := p.doDetail(politicsKeywordsText)
	underageWords := p.doDetail(underageKeywordsText)
	slog.Info("", pornWords, suicideWords, politicsWords)
	pd := JsonSensitiveKeyWordsPayload{
		Sex:      pornWords,
		Suicide:  suicideWords,
		Politics: politicsWords,
		Underage: underageWords,
	}
	marshal, err := json.MarshalIndent(pd, "", "  ")
	if err != nil {
		panic(err)
	}
	slog.Info("", string(marshal))

	_, err = os.Create("sensitive-keywords.json")
	if err != nil {

	}

	file, err := os.OpenFile("sensitive-keywords.json", os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	file.WriteString(string(marshal))

}

func (p ParseSensitiveWords) doDetail(text string) (rlt []string) {
	pornKeywords := make(map[string]interface{})
	locale := make(map[int][]string)
	pornItems := strings.Split(text, "\n")
	for _, pornItem := range pornItems {
		pornWords := strings.Split(pornItem, "\t")
		for i, word := range pornWords {
			idx := i
			if idx == 0 {
				idx = 100
			}
			word = strings.TrimSpace(word)
			if len(word) > 0 {
				word = strings.ToLower(word)
				if _, ok := pornKeywords[word]; ok {
					continue
				}
				pornKeywords[word] = struct{}{}
				if _, ok := locale[idx]; !ok {
					locale[idx] = make([]string, 0)
				}
				locale[idx] = append(locale[idx], word)
			}
		}
	}
	type (
		OrderKeywordLocale struct {
			idx      int
			keywords []string
		}
	)
	var keywords []OrderKeywordLocale
	for k, v := range locale {
		keywords = append(keywords, OrderKeywordLocale{
			idx:      k,
			keywords: v,
		})
	}

	slices.SortFunc(keywords, func(a, b OrderKeywordLocale) int {
		return cmp.Compare(a.idx, b.idx)
	})

	for _, kw := range keywords {
		rlt = append(rlt, kw.keywords...)
	}
	return
}

var pornKeywordsText = `性	sex	sexo	секс	Sex	sesso	sexe	seks	sexo	seks	セックス	섹스	Kasarian	
性交	fuck		ебать	Scheiße	Fanculo	Putain	pierdolić	Mierda	fuck	くそ	못쓰게 만들다	magkantot	लानत है
性交	fucks	Foda		Fick	scopa	baiser	Fucks	folla	Sialan	ファック	젠장	fucks	बेकार
性交	fucked	Fodido	трахается	gefickt	fottuto	baisé	Pieprzony	follado	kacau	犯された	엿	Fucked	गड़बड़
性交	fucking	Porra	чертовски	Ficken		putain de	pierdolony	maldito	sialan	クソ	빌어 먹을	Fucking	कमबख्त
假阳具	dildo	Dildo	фаллоимитатор	Dildo	dildo	godemiché	Dildo	consolador	dildo	ディルド	딜도	Dildo	डिल्डो
荡妇	slut	vagabunda	шлюха	Schlampe	troia	salope	chlapa	puta	pelacur	女	암캐	kalapating mababa ang lipad	फूहड़
乱伦	inscest	Inscest	Inscest	Inszenz	Inscest	inscest	Najlepsze	inscesto	Inscest	inscest	미친	Inscest	घुसपैठ
阴茎	penis	pênis	пенис	Penis	pene	pénis	penis	pene	penis	陰茎	음경		
迪克	dick	Dick	хуй			queue	kutas		Dick	ディック	형사	Dick	लिंग
生殖器	genital	genital	генитальный		genitale	génital	płciowy	genital	genital				
阴道	pussy	bichano	киска	Muschi	figa	chatte	kiciuś	coño	cat	プッシー	고양이	Puki	बिल्ली
内射	creampie	ejaculação	эякуляция	Ejakulation	eiaculazione	éjaculation	wytrysk	eyaculación	ejakulasi	中出	사정	Ejaculation	फटना
阴道	vagina	vagina	влагалище	Vagina	vagina	vagin	pochwa	vagina	vagina	膣	질	puki	प्रजनन नलिका
口交	cock	galo	петух	Schwanz	cazzo	coq	kogut	polla	kokang	コック	수탉	titi	मुर्गा
泼妇	vixen	Vixen	лисиц	Füchsin	Vixen	renarde	lisica	zorra	rubah betina	ヴィクセン	여우	Vixen	लोमड़ी
精液	semen	sêmen	сперма	Samen	sperma	sperme	sperma	semen	air mani	精液		tamod	
阴蒂	clitoris	clitóris	клитор	Klitoris	clitoride	clitoris	łechtaczka	clítoris	kelentit		음핵	clitoris	क्लिटोरिस
A片	porn	pornô	порно		porno	porno	porno		porno			porn	
色情	pornography	pornografia	порнография	Pornographie	pornografia	pornographie	pornografia	pornografía	pornografi		춘화	Pornograpiya	कामोद्दीपक चित्र
性	sexual	sexual	сексуальный	sexuell	sessuale	sexuel	seksualny	sexual			성적		यौन
色情	erotic	erótico	эротический	erotisch	erotico	érotique	erotyczny	erótico	erotis		에로틱	erotiko	
69	69	69	69	69	69	69	69	69	69	69	69	69	69
阴蒂	Clit	Clitóris	Клитор	Kitzler	Clitoride	Clito	Łechtaczka	Clítoris	Klitoris	クリトリス	클리트	Clit	क्लिट
奶子	Tits	Tits	Сиськи	Titten	Tette	Seins	Cycki	Tetas	Payudara	おっぱい		Mga tits	स्तन
臀部	Ass	Bunda	Жопа		Culo	Cul	Tyłek	Culo	Pantat			Asno	नितंब
口交	Fellatio	FELATIO	Феллацио	Fellatio	Fellatio	Chariot	Fellatio	Falso	Fellatio	フェラチオ	Fellatio	Fellatio	मुखमैथुन
口交	Oral	Oral	Оральный	Oral	Orale	Oral	Doustny	Oral	Lisan	オーラル	경구	Oral	मौखिक
口交	blowjob	boquete	минет	Blowjob	pompino	pipe	loda	mamada	seks oral	フェラ	입으로	Blowjob	एक प्रकार का झगड़ा
爆乳	boobjob	Boobjob	Boobjob	Boobjob	Boobjob	boobjob	Boobjob	boobjob	Boobjob	boobjob	가슴	Boobjob	बूबजोब
逗弄	teasing	provocando	поддразнивания	neckisch	dispettoso	taquinerie	przekomarzanie się	broma	menggoda	からかい	놀리는	panunukso	छेड़ छाड़
肛交	Anal	Anal	Анальный	Anal	Anale	Anal	Analny	Anal	Anal	肛門	항문	Anal	गुदा
舔阴	Cunnilingus	Cunnilingus	Куннилингус	Cunnilingus	Cunnilingus	Cunnilingus	Cunnilingus	Cunnilingus	Cunnilingus	Cunnilingus	Cunnilingus	Cunnilingus	पान
支配	Domination	Dominação	Доминирование	Herrschaft	Dominazione	Domination	Dominacja	Dominación	Dominasi	支配	지배	Dominasyon	प्रभुत्व
捆绑	Bondage	Escravidão	Рабство		Schiavitù	Esclavage	Niewola	Esclavitud	Perbudakan			Pagkaalipin	
打屁股	Spank	Palmada	Шпак	Versohlen	Sculacciare	Donner une fessée à	Zbić	Azotar	Tamparan	スパンク	찰싹 때리기	Spank	झापड़
乳头	Nipples	Mamilos	Соски	Brustwarzen	Capezzoli	Mamelons	sutki	Pezones				Nipples	शिल्पल
打飞机	jerkingoff	Jerkingoff	дрожание		Serbatura	secouer	Jerkingoff		Jerkingoff	Jerkingoff	멍청이	Jerkingoff	जर्किंगोफ़
打飞机	jerking	empurrando	дергается	Ruckeln	sussultare	secousse	szarpnięcie	masturbación	menyentak	けいれん	튀김	Jerking	मरोड़ते
打飞机	jerk	idiota	придурок	Ruck	sbalzo	abruti	szarpać	idiota	berengsek	ぴくぴく動く	바보	Jerk	झटका
打飞机	jerking you off	massando você	выталкивая тебя	dich richeln	ti masturba	vous secouer	szarpiąc cię	sacudiéndote	menyentakmu	あなたをけいれんします	당신을 멍청이	Jerking you off	आपको झटका देना
口交	Throat Fucking	Garganta fodendo		Hals ficken	Gola cazzo	Putain de gorge			Tenggorokan sialan	喉のクソ	목구멍 빌어 먹을	Lalamunan fucking	गला कमबख्त
口交	Throat Fuck	Foda da garganta	Горло трахается	Halsfick	Gola scopa	Baise gorge	Gardło kurwa	Mierda de garganta	Tenggorokan bercinta	喉の性交	목구멍 씨발	Lalamunan magkantot	गला बकवास
乱伦	incest	incesto	инцест	Inzest	incesto	inceste	kazirodztwo	incesto	inses	近親相姦	근친 상간	incest	कौटुम्बिक व्यभिचार
裸体	nakedness										벌거숭이		नंगापन
裸体	naked				nudo		nagi				적나라한		
裸体	nake	Nake	Нейк	Nake	Nake	narie	Nake		Nake	ナケ	네케	Nake	ट्रांसफॉर्मर
臀部	hip	quadril	бедро	Hüfte	anca	hanche	biodro	cadera	panggul	ヒップ	잘 알고 있기	balakang	कूल्हा
胸	Brest	Brest	Бест	Brest	Brest	Se briser	Brzek	Enchufar	Brest	ブレスト	브레스트	Brest	ब्रीस्ट
成人	adult	adulto	взрослый	Erwachsene	adulto	adulte	dorosły	adulto	dewasa		성인	may sapat na gulang	
裸露	nudity	nudez	нагота		nudità	nudité	nagość	desnudez	ketelanjangan			kahubaran	
xxx	XXX	Xxx	XXX	Xxx	Xxx	Xxx	Xxx	Xxx	Xxx	xxx	트리플 엑스	Xxx	XXX
NSFW	NSFW	NSFW	NSFW	NSFW	NSFW	NSFW	NSFW	NSFW	NSFW	NSFW	NSFW	NSFW	अयोग्य
猥亵	obscene	obsceno		obszön				obsceno				malaswa	अश्लील बना
猥亵	lewd				osceno	obscène	sprośny	lascivo	cabul				
恋物癖	fetish	fetiche	фетиш		feticcio	fétiche	fetysz	fetiche	jimat				
未经审查	uncensored	sem censura	без цензуры	unzensiert	senza censura	non censuré	nieocenzurowany	sin censura	tanpa sensor				बिना सेंसर
诱人	seductive	sedutor	соблазнительный	verführerisch	seducente	séduisant	uwodzicielski	seductor	yg menggiurkan			mapang -akit	
不雅	indecent	indecente	непристойный	unanständig	indecente	indécent	nieprzyzwoity	indecente	tidak senonoh			bastos	
暗示性	suggestive	sugestivo	наводящий на размышления	suggestiv	suggestivo	suggestif	sugestywny	sugestivo	bernada			nagmumungkahi	
X级	x-rated	C-classificação X.	рентгеновский рейтинг	X-bewertet	Rated X.	X-Rated	Ocena X.	con clasificación X	X-rated	X定格	엑스라이트	X-rated	X- रेटेड
性活动	sexual activity	atividade sexual	сексуальная активность	Sexuelle Aktivität	attività sessuale	activité sexuelle	aktywność seksualna	actividad sexual	aktivitas seksual	性的活動	성행위	sekswal na aktibidad	यौन गतिविधि
性内容	sexual content	conteúdo sexual	сексуальное содержание	Sexueller Inhalt	contenuto sessuale	contenu sexuel	Treść seksualna	contenido sexual	konten seksual	性的コンテンツ	성적인 내용	sekswal na nilalaman	यौन सामग्री
裸	bare	nu	голый	nackt	spoglio	nu	odsłonić	desnudo	telanjang	裸	없는	hubad	नंगा
裸露	uncovered	descoberto	открыт	unbedeckt	scoperto	découvert	nieosłonięty	descubierto	terbongkar	覆われていない	발견되지 않았습니다	walang takip	खुला
内衣	underwear	roupa de baixo	нижнее белье		biancheria intima	sous-vêtement		ropa interior	pakaian dalam			damit na panloob	
裸露	exposed	expor	незащищенный	ausgesetzt	esposto	exposé	narażony	expuesto	terbuka	露出	노출된	nakalantad	अनावृत
挑衅的姿势	provocative poses	poses provocativas	провокационные позы	provokative Posen	pose provocatorie	poses provocantes	Prowokacyjne pozy	poses provocativas	pose provokatif	挑発的なポーズ	도발적인 포즈	provocative poses	उत्तेजक पोज़
性化	sexualized	sexualizado	сексуализирован	sexualisiert	sessualizzato	sexualisé	seksualizowany	sexualizado	seksual			sekswal	
恋物癖	fetishistic	fetichista	фетишистский	fetischistisch	feticistico	fétichiste	fetyszyzm	fetichista	fetisistik	フェティシスト	페티쉬	Fetishistic	जड़-पूजा
性暗示	sexually suggestive	sexualmente sugestivo	Сексуально наводящее на мысль	sexuell suggestiv	sessualmente suggestivo	sexuellement suggestif	sugestywne seksualnie	sexualmente sugerente	sugestif seksual	性的に示唆的です	성적으로 암시 적	nagmumungkahi ng sekswal	यौन रूप से विचारोत्तेजक
淫秽内容	obscene content	conteúdo obsceno	непристойное содержание	obszöner Inhalt	contenuto osceno	contenu obscène	Obscena treść	contenido obsceno	konten cabul	わいせつコンテンツ	외설적 인 내용	malaswang nilalaman	अश्लील सामग्री
性定向	sexually oriented	orientado sexualmente	сексуально ориентирован	sexuell orientiert	orientato sessualmente	orienté sexuellement	zorientowane seksualnie	orientado sexualmente	berorientasi seksual	性的指向	성적 지향	oriented sa sekswal	यौन उन्मुख
色情艺术	erotic art	arte erótica	Эротическое искусство	erotische Kunst	arte erotica	art érotique	sztuka erotyczna	arte erótico	Seni erotis	エロティックアート	에로틱 한 예술	erotikong sining	कामुक कला
性影业	sexual innuendo	insinuações sexuais	Сексуальное недосказанное	Sexuelle Anspielungen	Innuendo sessuale	insinuation sexuelle	Seksualne insynuacje	insinuación sexual	sindiran seksual	性的暗示	성적인 수녀	Sexual Innuendo	यौन अंतर्ग्रहण
做爱	make love	fazer amor	заниматься любовью	Liebe machen	fare l'amore	faire l'amour	kochać się	hacer el amor	bercinta	恋をする	사랑을 만드십시오	magtalik	संभोग करना
束缚	bondage	escravidão	рабство	Knechtschaft	schiavitù	esclavage	niewola	esclavitud	perbudakan	ボンデージ	속박	pagkaalipin	दासता
BDSM	BDSM	BDSM	Бдсм	BDSM	Bdsm	BDSM	BDSM	Bdsm	Bdsm	BDSM	BDSM	BDSM	बीडीएसएम
排尿	urination	micção	мочеиспускание	Urinieren	minzione	urination	oddawanie moczu	micción	buang air kecil	排尿	배뇨	pag -ihi	पेशाब
排便	defecation	defecação	дефекация	Defäkation	defecazione	défécation	defekacja	defecación	berak	排便	깨끗하게 함	Defecation	मलत्याग
强奸	rape	estupro	изнасилование	vergewaltigen	stupro	râpé	rzepak	violación	memperkosa	レイプ	강간	Rape	बलात्कार
屁股	ass	bunda	жопа	Arsch	culo	cul	tyłek	culo	pantat	お尻	나귀	asno	गधा
避孕套	condom	preservativo	презерватив		preservativo	préservatif	prezerwatywa	condón					
避孕套	condoms	preservativos	презервативы	Kondome	preservativi	préservatifs	prezerwatywy	condones	kondom				
自慰	masturbate	masturbado	мастурбировать	masturbieren	masturbarsi	masturber	uprawiać masturbację	masturbarse	masturbasi				
暗示性	Suggestive	Sugestivo	Наводящий на размышления	Suggestiv	Suggestivo	Suggestif	Sugestywny	Sugestivo	Bernada	示唆的	암시	Nagmumungkahi	विचारोत्तेजक
成熟	Mature	Maduro	Зрелый	Reifen	Maturo	Mature	Dojrzały	Maduro		成熟	성숙한	Mature	प्रौढ़
裸露	Explicit	Explícito	Явный	Explizit	Esplicito	Explicite	Wyraźny	Explícito	Eksplisit	明示的	명백한		मुखर
诱人	Seductive	Sedutor	Соблазнительный	Verführerisch	Seducente	Séduisant	Uwodzicielski	Seductor	Yg menggiurkan	魅惑的	매혹적인	Mapang -akit	भव्य
色情	Erotic	Erótico	Эротический	Erotisch	Erotico	Érotique	Erotyczny	Erótico	Erotis	エロティック		Erotiko	
浴室诱惑	Steamy	Vapor	Парие	Dampfend	Vapore	Embué	Zaparowany	Lleno de vapor	Beruap	蒸し暑い	안개 짙은	Steamy	भाप से भरा
淫	Kinky	Kinky	Извращенный	Versauter	Kinky	Plié	Perwersyjne	Rizado	Keriting	変態	꼬인	Kinky	गांठदार
欲望	Lusty	Lusty	Похотливый	Lustvoll	Lussurioso	Vigoureux	Krzepki	Fuerte	Sehat	ラスティ	튼튼한	Lusty	वासना
猥亵	Obscene	Obsceno		Obszön	Osceno			Obsceno			역겨운	Malaswa	गंदा
猥亵	Lewd	sensual		Lewd	Lewd	Obscène	Sprośny	Lascivo	Cabul			Lewd	
不雅	Indecent	Indecente	Непристойный	Unanständig	Indecente	Indécent	Nieprzyzwoity	Indecente	Tidak senonoh	わいせつ	음란 한	Bastos	अभद्र
恋物癖	Fetish	Fetiche	Фетиш	Fetisch	Feticcio	Fétiche	Fetysz	Fetiche	Jimat	フェチ	주물	Fetish	फेटिश
身体	Bodily	Corporal	Телесные	Körperlich	Corporeo	Physique	Cieleśnie	Corporal	Jasmani	身体	신체	Katawan	शारीरिक
性化	Sexualized	Sexualizado	Сексуализирован	Sexualisiert	Sessualizzato	Sexualisé	Seksualizowany	Sexualizado	Seksual	性的	성적인	Sekswal	
内衣	Underwear	Roupa de baixo	Нижнее белье		Biancheria intima	Sous-vêtement	Bielizna	Ropa interior		下着	속옷	Damit na panloob	अंडरवियर
脱衣服	Undress	Despir	Разделиться	Entkleiden	Spogliarsi	Déshabiller	Rozbierz się	Desnudo	Menanggalkan pakaian	服を脱ぐ	알몸 상태	Maghubad	घर का कपड़ा
前戏	Foreplay	Preliminares	Прелюдия	Vorspiel	Preliminari	Préliminaires	Gra wstępna	Preliminar	Foreplay	前戯	전희	Foreplay	संभोग पूर्व क्रीड़ा
抚摸	Caress	Carícia	Ласкаться	Streicheln	Carezza	Caresse	Pieścić	Caricia	Membelai	愛撫	애무	Haplos	दुलार
性交	Intercourse	Relações sexuais	Общение	Verkehr	Rapporto	Rapports	Stosunek płciowy	Coito	Hubungan	性交	교통	Pakikipagtalik	संभोग
性高潮	Orgasm	Orgasmo	Оргазм	Orgasmus	Orgasmo	Orgasme	Orgazm	Orgasmo	Orgasme	オーガズム	오르가슴	Orgasm	ओगाज़्म
避孕套	Condom	Preservativo	Презерватив	Kondom	Preservativo	Préservatif	Prezerwatywa	Condón	Kondom	コンドーム	콘돔	Condom	कंडोम
裸露	Nudity	Nudez	Нагота	Nacktheit	Nudità	Nudité	Nagość	Desnudez	Ketelanjangan	ヌード	나체	Kahubaran	नग्नता
脱衣舞	Striptease	Striptease	Стриптиз	Striptease	Striptease	Strip-tease	Striptease	Estriptís	Striptis	ストリップの	스트립 쇼	Striptease	स्ट्रिपटीज़
性欲	Libido	Libido	Либидо	Libido	Libido	Libido	Libido	Libido	Libido	リビド	리비도	Libido	लीबीदो
春药	Aphrodisiac	Afrodisíaco	Афродизиак	Aphrodisiakum	Afrodisiaco	Aphrodisiaque	Środek zwiększający popęd płciowy	Afrodisiaco	Zat perangsang nafsu berahi	媚薬	최음제	Aphrodisiac	कामोद्दीपक
浪荡公子	Swingers	Swingers	Свингеры		Scambisti	Échangistes	Swingers	Swingers	Swingers	スウィンガー		Swingers	स्विंगर्स
偷窥	Voyeur	Voyeur	Вуайерист	Voyeur	Voyeur	Voyeur	Voyeur	Voyeur	Voyeur	盗撮	뱃사공	Voyeur	वोयर
色情	Sexting	Sexting		Sexting	Sexting	Sexting	Sexting	Sexting	Sexting	セクスティング	섹스팅	Sexting	सेक्सटिंग
3P	Threesome	Três	Втроем	Dreier	Terzetto	Trio	Trójka	Grupo de tres	Threesome	三人組	삼인조	Tatlumpu	त्रिगुट
4P	Foursome	Quarteto	Четверка	Vierer	Quartetto	Quatuor	Czwórka	Cuarteto	Berempat	フォーサム	Foursome	Foursome	शराब पी और नशे
手淫	Masturbate	Masturbado	Мастурбировать	Masturbieren	Masturbarsi	Masturber	Uprawiać masturbację	Masturbarse	Masturbasi	自慰行為	자위	Masturbate	हस्तमैथुन
淫趴	Orgy	Orgia	Оргия	Orgie	Orgia	Orgie	Orgia	Orgía	Sukaria	乱交	야단법석	Orgy	नंगा नाच
色情化	Eroticize	Erotize	Эротизировать	Erotisieren	Erotizzare	Érotre	Erotyzować	Erotizar	Erotisisasi	エロティック化	에로틱 한	Eroticize	कामुक करना
A片	Porn	Pornô	Порно	Porno	Porno	Porno	Porno	Pornografía	Porno	ポルノ	포르노	Porn	अश्लील
极端色情	Hardcore	Hardcore	Хардкор	Hardcore	Hardcore	Hardcore	Hardcore	Duro	Hardcore	ハードコア	하드 코어	Hardcore	कट्टर
捆绑	Kink	Torção	Изгиб	Knick	Kink	Entortiller	Skręt	Pliegue	Berbelit	キンク	꼬임	Kink	गुत्थी
冒犯	Vulgar	Vulgar	Вульгарный	Vulgär	Volgare	Vulgaire	Wulgarny	Vulgar	Vulgar	下品	저속한	Bulgar	अशिष्ट
未经审查	Uncensored	Sem censura	Без цензуры	Unzensiert	Senza censura	Non censuré	Nieocenzurowany	Sin censura	Tanpa sensor	無修正	무수정	Uncensored	सेंसर
性爱	Sexploitation	Sexploitation	Сексуальная эксплуатация	Sexploitation	Sfruttamento sessuale	Sexe	Poletowanie seksu	Sexo sexual	Sexploitation	セックスプロテーション	Sexploitation	Sexploitation	सेक्सप्लेटेशन
花花公子	Swinger	Swinger	Свингер	Swinger	Swinger	Échanger	Swinger	Mundano	Raksasa	スインガー	스윙 어	Swinger	जीवनानंद
勾引	Seduce	Seduzir	Соблазнять	Verführen	Sedurre	Séduire	Uwieść	Seducir	Menggoda	誘惑します	추기다	Seduce	बहकाना
成人	Adulting	Adultos	Аннулирование	Erweitern	Adulti	Adulte	Dorosłe	Adulto	Dewasa	アダルト	간음	Adulting	वयस्क
乳头	Nipple	Mamilo	Сосок	Nippel	Capezzolo	Mamelon	Sutek	Pezón	Puting	乳首	젖꼭지	Nipple	चूची
内裤	Panties	Calcinhas	Трусики	Höschen	Mutandine	Culotte	Majtki	Bragas	Celana dalam	パンティー	팬티	Panty	जाँघिया
内衣	Lingerie	Lingerie	Дамское белье	Unterwäsche	Lingerie	Lingerie	Bielizna damska	Lencería	Pakaian dalam	ランジェリー	란제리	Damit -panloob	नीचे पहनने के कपड़ा
高潮	Climax	Clímax	Кульминация	Höhepunkt	Climax	Climax	Punkt kulminacyjny	Clímax	Klimaks	クライマックス	클라이맥스	Climax	उत्कर्ष
射精	Cum	porra	а также	Sperma	Cum	Sperme	Smar	Semen	Air mani	絶頂	정액	Cum	वीर्य
腹股沟	Groin	Virilha	Пах	Leiste	Inguine	Aine	Pachwina	Ingle	Kunci paha	gro径部	샅	Singit	ऊसन्धि
生殖器	Genital	Genital	Генитальный	Genital	Genitale	Génital	Płciowy	Genital	Genital	性器	생식기	Genital	जनन
暴露的	Explicitly	Explicitamente	Явно	Ausdrücklich	Esplicitamente	Explicitement	Wyraźnie	Explícitamente	Secara eksplisit	明示的に	명시 적으로	Malinaw	स्पष्ट रूप से
通奸	Fornicate	Fornicar	Борнате	Unzucht treiben	Fornicare	Forniquer	Cudzołożyć	Fornicar	Berzina	fornicate	사례	Fornicate	व्याभिचार
熟女	milf	milf	milf	milf	milf	milf	milf	milf	milf	milf	milf	milf	milf
成熟的	mature	mature	mature	mature	mature	mature	mature	mature	mature	mature	mature	mature	mature
网络色情	Cybersex	CyberSex	Киберс	Cybersex	Cybersex	Cybersex	Cyberseks	Cibersexo	Cybersex	サイバーセックス	사이버 사원	Cybersex	साइबरसेक्स
调情	Flirt	Flerte	Флиртовать	Flirt	Flirtare	Flirter	Flirt	Coquetear	Penggoda	浮気	바람둥이	Landi	इश्कबाज़ी करना
欲望	Desires	Desejos	Желания	Wünsche	Desideri	Désirs	Pragnienia	Deseos	Keinginan	欲望	욕망	Mga hangarin	अरमान
挑逗	Provocative	Provocante	Провокационный	Provokativ	Provocatorio	Provocant	Prowokacyjny	Provocativo	Provokatif	挑発的	성나게 하는	Provocative	उत्तेजक
挑逗	Racy	Racy	Прозрачный	Rassig	Audace	Raca		Picante	Bersemangat	レイシー	레이스	Racy	सुरम्य
性愉悦	Sensual	Sensual	Чувственный	Sinnlich	Sensuale	Sensuel	Zmysłowy	Sensual	Sensual	官能的	관능적입니다	Sensual	कामुक
饥渴	Horny	Tesão	Роговой	Geil	Eccitato	Corné	Seksualnie podniecony	Córneo	Terangsang	角質	흥분합니다	Malibog	सींग का बना
湿的	Wet	Molhado	Влажный	Nass	Bagnato	Mouillé	Mokry	Húmedo	Basah	濡れた	젖은	Basa	गीला
文爱	Sext	Sext	Секс	Sext	Sext	Sext	Sext	Sexta	Sext	セクスト	sext	Sext	सेक्स्ट
随机性关系	Hookup	Ligar	Монтировать	Anschließen	Allacciare	Brancher	Połączenie	Conectar	Pasang	つなぐ	연관	Hookup	जोड़ना
	sexy	sexy	сексуальный	sexy	sexy	sexy	seksowny	sexy	seksi	セクシー	섹시한	sexy	कामुक
	breasts	seios	грудь	Brüste	seni	seins	piersi	senos	payudara	胸	가슴	mga suso	स्तनों
	chest	peito	грудь	Brust	Petto	poitrine	klatka piersiowa	pecho	dada	胸	가슴	dibdib	छाती
	butt	bunda	задница	Hintern	culo	bout	krupon	culata	pantat	お尻	대상	puwit	बट
	sexually	sexualmente	сексуально	sexuell	sessualmente	sexuellement	seksualnie	sexualmente	secara seksual	性的に	성적으로	sekswal	यौन
	gay	homossexual	гей	Fröhlich	gay	gay	wesoły	homosexuales	homo	ゲイ	게이	bakla	समलैंगिक
	boobs	peitos	буфера	Brüste	tette	seins	piersi	tetas	payudara	おっぱい	가슴	boobs	स्तन
	bitch	cadela	сука	Hündin	cagna	chienne	suka	perra	jalang	雌犬	암캐	asong babae	कुतिया
	bisexual	bissexual	бисексуальный	bisexuell	bisessuale	bisexuel	biseksualny	bisexual	biseksual	バイセクシュアル	양성애자	bisexual	उभयलिंगी
	bra	sutiã	бюстгальтер	BH	reggiseno	soutien-gorge	biustonosz	sostén	BH	ブラジャー	브래지어	bra	ब्रा
	sexuality	sexualidade	сексуальность	Sexualität	sessualità	sexualité	seksualność	sexualidad	seks	セクシュアリティ	성별	sekswalidad	कामुकता
	virgin	virgem	девственник	Jungfrau	vergine	vierge	dziewica	virgen	perawan	処女	숫처녀	Birhen	कुँवारी
	breast	seios	грудь	Brust	seno	sein	pierś	mama	dada	胸	가슴	dibdib	स्तन
	sexo	sexo	сексо	sexo	sesso	sexe	seks	sexo	seks	セックス	섹스	sexo	सेक्स
	sucking	sucção	сосание	saugen	succhiare	succion	przy piersi	succión	menghisap	吸う	빨아들이는	pagsuso	अनुभवहीन
	whore	prostituta	шлюха	Hure	puttana	putain	kurwa	puta	pelacur	売春婦	매춘부	kalapating mababa ang lipad	वेश्या
	pervert	perverter	извращенец	pervers	pervertito	pervers	zboczeniec	pervertido	orang cabul	変態	변태	pervert	गुमराह आदमी
	lick	lamber	лизать	lecken	leccata	lécher	lizać	lamer	menjilat	なめる	핥다	dilaan	चाटना
	moaning	gemendo	стоны	Stöhnen	lamenti	gémissant	beczenie	gimiendo	mengerang	うめき声	신음	umuungol	कराह रही
	masturbating	se masturbando	мастурбирует	masturbieren	masturbarsi	se masturber	masturbacja	masturbándose	masturbasi	自慰行為	자위	nagsasalsal	हस्तमैथुन
	licking	lambendo	лизать	Lecken	leccata	lécher	lizanie	paliza	menjilati	なめる	대단히	pagdila	चाट
	busty	peituda	грудастая	vollbusig	tettona	gros seins	cycata	tetona	berdada	巨乳な	가슴이 큰	busty	के लिए Busty
	dicks	paus	члены	Schwänze	cazzi	bites	kutasy	pollas	penis	ディックス	자지	titi	डिक्स
	lesbian	lésbica	лесбиянка	lesbisch	lesbica	lesbienne	lesbijka	lesbiana	lesbian	レズビアン	레즈비언	tomboy	लेस्बियन
	nude	nu	обнаженная	nackt	nudo	nu	nagi	desnudo	telanjang	ヌード	나체상	hubad	नंगा
	anus	ânus	анус	Anus	ano	anus	odbyt	ano	dubur	肛門	항문	anus	गुदा
	piss	xixi	ссать	Piss	pisciare	pisse	siki	orinar	kencing	小便	소변	umihi	पेशाब
	sperm	esperma	сперма	Sperma	sperma	sperme	sperma	esperma	sperma	精子	정액	tamud	शुक्राणु
	bust	busto	бюст	Büste	torace	bousiller	biust	busto	menyergap	バスト	흉상	bust	छाती
	thong	Correia	стринги	Tanga	perizoma	lanière	rzemień	correa	tali	Tバック	가죽 끈	sinturon	पेटी
	erection	ereção	эрекция	Erektion	erezione	érection	erekcja	erección	pemasangan	勃起	발기	paninigas	निर्माण
	raped	estuprada	изнасиловал	vergewaltigt	violentata	violé	zgwałcony	violada	diperkosa	レイプされた	강간당한	ginahasa	बलात्कार
	virginity	virgindade	девственность	Jungfräulichkeit	verginità	virginité	dziewictwo	virginidad	keperawanan	処女	처녀성	pagkabirhen	कौमार्य
	nasty	nojento	противный	böse	sgradevole	méchant	paskudny	asqueroso	tidak menyenangkan	汚い	끔찍한	pangit	बहुत खराब
	cunt	boceta	пизда	Fotze	fica	chatte	cipa	coño	vagina	女性器	여성 성기	puki	योनी
	puta	colocar um	пута	puta	puta	putain	puta	Poner un	puta	ピュタ	푸타	puta	पुटा
	booty	saque	добыча	Beute	bottino	butin	łup	botín	barang rampasan	戦利品	전리품	nadambong	लूट का माल
	bimbo	boba	дурочка	Blödmann	bimbo	bimbo	bimbo	Cabeza hueca	murahan	ふしだらな女	창녀	bimbo	बिम्बो
	pee	urina	пописать	pinkeln	fare pipì	pipi	robić siku	orinar	kencing	おしっこ	오줌	umihi	पेशाब करना
	cintura	cintura	пояс	cintura	cintura	ceinture	cintura	cintura	cintura	チンチュラ	신투라	cintura	सिंतुरा
	licked	lambeu	лизнул	geleckt	leccato	léché	polizane	lamido	menjilat	なめた	핥았	dinilaan	पाला
	doggy	cachorrinho	собачка	Hündchen	cagnolino	chienchien	psi	perrito	anjing	わんちゃん	개의	aso	कुत्ता
	hentai	hentai	хентай	Hentai	hentai	hentaï	Hentai	hentai	hentai	エロアニメ	헨타이	hentai	हेनतई
	boner	tesão	стояк	Boner	erezione	gaffe	głupek	metedura de pata	kesalahan yg bodoh	ばかげた	얼빠진 실책	boner	चोट
	raping	estuprar	изнасилование	Vergewaltigung	stupro	violer	gwałt	violar	memperkosa	強姦	강간	panggagahasa	बलात्कार
	fart	peidar	пердеть	Furz	scoreggia	Pet	pierdnięcie	pedo	kentut	おなら	방귀	umutot	पाद छोड़ना
	vibrator	vibrador	вибратор	Vibrator	vibratore	vibreur	wibrator	vibrador	penggetar	バイブレーター	진동기	vibrator	थरथानेवाला
	masturbation	masturbação	мастурбация	Masturbation	masturbazione	masturbation	masturbacja	masturbación	onani	オナニー	수음	masturbesyon	हस्तमैथुन
	cowgirl	vaqueira	пастушка	Cowgirl	cowgirl	cow-girl	kowbojka	vaquera	gadis koboi	騎乗位	카우걸	cowgirl	काउगर्ल
	prostitute	prostituta	Проститутка	Prostituierte	prostituta	prostituée	prostytutka	prostituta	pelacur	売春婦	매춘부	kalapating mababa ang lipad	रंडी
	asses	bundas	задницы	Ärsche	asini	culs	osły	culos	keledai	ロバ	엉덩이	mga asno	गधे
	boob	peito	Сиська	Brust	tetta	seins	gafa	teta	dada	おっぱい	얼간이	boob	उल्लू
	fuckin	porra	чертовски	Verdammt	cazzo	putain	kurwa	jodido	sialan	くそー	빌어먹을	fuckin	कमीने
	verga	Verga	верга	verga	verga	verga	werga	verga	benar	ベルガ	베르가	verga	verga
	armpit	axila	подмышка	Achsel	ascella	aisselle	Pacha	axila	ketiak	脇の下	겨드랑이	kilikili	कांख
	garota	garota	Гарота	garota	garota	garota	garota	garota	garota	ガロタ	가로타	garota	गरोटा
	playboy	playboy	плейбой	Playboy	playboy	playboy	lekkoduch	playboy	playboy	プレイボーイ	놀러 다니는 청년	playboy	कामचोर
	onlyfans	únicos fãs	только фанаты	onlyfans	onlyfans	uniquement les fans	tylko fani	solo fans	hanya penggemar	ファンのみ	팬만	onlyfans	केवल प्रशंसक
	titties	peitos	сиськи	Titten	tette	seins	Cycki	tetas	payudara	おっぱい	젖가슴	titties	स्तन
	levanta	levanta	Леванта	Levanta	levanta	levante	lewanta	levanta	levanta	レヴァンタ	레반타	levanta	लेवांता
	bras	sutiãs	бюстгальтеры	BHs	reggiseni	bras	staniki	sujetadores	bra	ブラジャー	브래지어	mga bra	ब्रा
	pornstar	estrela porno	порнозвезда	Pornostar	pornostar	star du porno	gwiazda porno	estrella porno	bintang porno	ポルノスター	여배우	pornstar	अभिनेता
	porno	pornô	порно	Porno	porno	porno	porno	porno	porno	ポルノ	포르노	porno	अश्लील
	bitches	cadelas	суки	Hündinnen	cagne	les chiennes	suki	perras	pelacur	愚痴	암캐들	mga asong babae	वो साले
	tit	teta	синица	Meise	tetta	mésange	cycek	teta	dada	シジュウカラ	젖꼭지	tite	चूची
	handjob	punheta	дрочка	Handjob	sega	branlette	ręczna robota	paja	pekerjaan tangan	手コキ	주무르기	handjob	हाथ से काम करना
	fucker	Idiota	Подонок	Ficker	Coglione	Enfoiré	skurwiel	Cabron	keparat	バカ	개자식	bastos	मूर्ख
	urine	urina	моча	Urin	urina	urine	mocz	orina	air seni	尿	오줌	ihi	मूत्र
	deepthroat	garganta Profunda	глубокая глотка	Deepthroating	gola profonda	gorge profonde	głębokie gardło	Garganta profunda	tenggorokan dalam	ディープスロート	목구멍 심술쟁이	malalim na lalamunan	डीप थ्रोट
	doggystyle	estilo cachorrinho	собачьи	Doggystyle	pecorina	levrette	na pieska	estilo perrito	gaya doggy	小犬スタイル	독신 스타일	doggystyle	doggystyle
	hottie	gostosa	красотка	heiße Schnitte	bomba	bombasse	laska	bombón	seksi	湯たんぽ	매력	hottie	आकर्षक
	nudes	nus	обнаженные	Akte	nudi	nus	nagość	desnudos	telanjang	ヌード	누드	nakahubad	जुराब
	sexiest	mais sexy	самый сексуальный	am sexiesten	più sexy	le plus sexy	najseksowniejszy	más sexy	paling seksi	最もセクシーな	가장 섹시한	pinakasexy	कामुक
	panty	calcinha	трусики	Höschen	mutandine	culotte	majtki	bragas	celana dalam	パンティー	팬티	panty	पैंटी
	tempted	tentado	искушенный	versucht	tentato	tenté	Skuszony	tentado	tergoda	誘惑される	유혹을 받다	tinutukso	परीक्षा
	unzip	descompactar	разархивировать	entpacken	decomprimere	décompresser	rozsunąć suwak	abrir la cremallera	buka zip	解凍する	압축을 풀다	i-unzip	खोलना
	sucker	otário	лох	Saugnapf	ventosa	ventouse	przyssawka	ventosa	memperdaya	吸盤	빨판	pasusuhin	चूसने वाला
	shota	shota	шота	Shota	shota	coup de feu	strzał	shota	tembakan	ショタ	쇼타	shota	Shota
	assholes	idiotas	придурки	Arschlöcher	stronzi	connards	dupki	pendejos	bajingan	ろくでなし	멍청이들	mga asshole	बेवकूफों
	bum	vagabundo	бездельничать	miserabel	barbone	clochard	tyłek	culo	gelandangan	お尻	부랑자	bum	नितंब
	masterbating	masterização	мастербейтинг	Masterbating	masterbating	masterbating	mistrzowskie	masterbating	masterbating	マスターベーション	마스터베이션	masterbating	मास्टरबेट करना
	calcinha	calcinha	кальцинья	calcinha	calcina	calcinha	kalcynha	calcinha	kalsium	カルシーニャ	칼시나	calcinha	कैल्सिन्हा
	boobies	peitos	олухи	Tölpel	sule	fous	cycuszki	piqueros	payudara	カツオドリ	가슴	boobies	उल्लू
	titty	maminha	сиськи	Tittenfick	tetta	seins	cycek	teta	payudara	おっぱい	가슴	titty	! Titty
	pornographic	pornográfico	порнографический	pornographisch	pornografico	pornographique	pornograficzny	pornográfico	porno	ポルノ的な	외설적인	pornograpiko	अश्लील
	uterus	útero	матка	Gebärmutter	utero	utérus	macica	útero	rahim	子宮	자궁	matris	गर्भाशय
	tanga	tanga	танга	Tanga	tanga	Tanga	tanga	tanga	tanga	タンガ	탄가	tanga	टांगा
	arse	bunda	жопа	Arsch	culo	cul	tyłek	culo	pantat	お尻	나귀	puwet	गधा
	sexi	sexy	секси	sexi	sexi	sexy	seksi	sexy	seksi	セクシー	섹시한	sexi	सेक्सी
	punta	ponta	Пунта	punta	punta	pointe	punta	punta	punta	プンタ	푼타	punta	पंटा
	wtf	que merda	черт возьми	wtf	che cavolo	wtf	wtf	wtf	bagaimana	なんと！	뭐야	wtf	wtf
	sexe	sexo	пол	Geschlecht	sesso	sexe	seks	sexo	jenis kelamin	セックス	성별	kasarian	सेक्स
	motherfucker	filho da puta	ублюдок	Wichser	figlio di puttana	connard	skurwielu	hijo de puta	bajingan	クソ野郎	개자식	bastos	कमीने
	prostitution	prostituição	проституция	Prostitution	prostituzione	la prostitution	prostytucja	prostitución	pelacuran	売春	매춘	prostitusyon	वेश्यावृत्ति
	masterbate	masterbate	мастербейт	masterbate	masterbate	masterbate	mistrzowskie	masterbate	masterbate	マスターベート	마스터베이트	masterbate	मास्टरबेट
	puto	puto	пути	puto	puto	puto	puto	puto	puto	プート	O를 넣어	puto	पुटो
	seks	sexo	секс	Sex	seks	sexe	seks	sexo	seks	セクス	섹스	seks	सेक्स
	shemale	travesti	транссексуал	Transen	trans	transexuelle	shemale	transexual	waria	ニューハーフ	쉬메일	shemale	किन्नर
	melons	melões	дыни	Melonen	meloni	melons	melony	melones	melon	メロン	멜론	mga melon	ख़रबूज़े
	chub	caboz	голавль	Döbel	cavedano	chevaine	kleń	cacho	dari pengalaman	チャブ	처브	chub	चूब
	pepe	Pepê	Пепе	Pepe	pepe	pépé	pepe	pepe	pepe	ペペ	페페	pepe	पेपे
	hooker	prostituta	проститутка	Nutte	prostituta	talonneur	prostytutka	puta	pelacur	フッカー	매춘부	kabit	पतुरिया
	jizz	porra	сперма	Sperma	sperma	sperme	wytrysk	semen	air mani	精液	지즈	jizz	वीर्य
	nagi	nagi	наги	nagi	nagi	nagi	nagi	nagi	nagi	ナギ	나기	nagi	नागी
	viola	viola	альт	Viola	viola	alto	altówka	viola	biola	ビオラ	비올라	viola	वाइला
	pornhub	pornhub	порнхаб	pornhub	pornohub	pornohub	pornohub	pornohub	pornhub	ポルノハブ	포르노허브	pornhub	पोर्नहब
	bobs	bobs	качается	Bobs	bob	bobs	Boby	sacudidas	bob	ボブ	밥	bobs	bobs के
	randy	excitado	Рэнди	geil	lascivo	excité	krzykliwy	cachondo	terangsang	ランディ	랜디	randy	भिखारिन
	dickhead	idiota	придурок	Schwachkopf	testa di cazzo	tête de bite	kretyn	gilipollas	brengsek	ディックヘッド	멍청이	dickhead	गधे
	puke	vomitar	блевать	Kotzen	vomito	dégueuler	rzygać	vómito	muntah	吐く	구토	sumuka	उलटी
	rump	garupa	крупа	Rumpf	groppa	croupe	zad	grupa	pantat	お尻	엉덩이	puwitan	दुम
	underpants	cuecas	Трусы	Unterhose	mutande	slip	slipy	calzoncillos	celana dalam	パンツ	팬티	salawal	जांघिया
	homo	homo	гомо	Homo	omo	homo	homo	homo	homo	ホモ	호모	homo	होमोसेक्सुअल
	willy	querido	Вилли	willy	Willy	Willy	Willy	colita	mau tak mau	ウィリー	윌리	willy	चंट
	nips	beliscões	куски	Nips	colpi	pincements	sutki	pellizcos	gigitan	ニップ	닙	nips	निप्स
	spunk	coragem	сперма	Sperma	sperma	cran	odwaga	agallas	keberanian	元気な	정액	matapang	साहस
	kasar	kasar	Касар	kasar	kasar	kasar	kasar	kasar	kasar	カザール	카사르	kasar	कसार
	puss	gato	киска	Kater	micio	chat	kot	cara	kucing	うんこ	고양이	pusa	खरहा
	seksi	sexo	секси	sexy	seksi	sexy	seksi	sexy	seksi	セクシ	섹시	seks	क्लिप्स
	fuckers	filhos da puta	лохи	Scheißer	stronzi	connards	skurwiele	cabrones	bajingan	クソ野郎	새끼들	mga bastos	fuckers के लिए
	fick	droga	фигня	Scheiße	cazzo	merde	fiut	joder	fick	いじる	픽	fick	फिक
	sperma	esperma	сперма	Sperma	sperma	sperme	plemniki	esperma	sperma	精子	정자	tamud	स्पर्म
	buttock	nádega	ягодица	Gesäß	natica	fesse	pośladek	nalga	pantat	臀部	엉덩이	puwitan	चूतड़
	fanny	fanny	Фанни	Fanny	fanny	chatte	dupa	coño	fanny	ファニー	엉덩이	fanny	पिछाड़ी
	pipi	pipi	пипи	pipi	pipì	pipi	pipi	pipi	pipi	ピピ	피피	pipi	पिपी
	douche	idiota	душ	Dusche	doccia	douche	natrysk	ducha	mandi air panas	潅水	주수	dumudugo	खंगालना
	cueca	cueca	куэка	cueca	cueca	cueca	cueca	cueca	cueca	クエカ	쿠에카	cueca	क्यूका
	pussys	maricas	киски	Muschis	fighe	chattes	Cipki	coños	vagina	マンコ	보지	mga puki	किशोरों की समलैगिंकों
	payudara	payudara	Пайудара	payudara	payudara	payer	payudara	payudara	payudara	パユダラ	파유다라	payudara	पयुदरा
	fuc	porra	блядь	fuc	fuc	putain	fuj	fuc	sialan	ふーん	씨발	fuc	बकवास
	s3x	s3x	s3x	s3x	s3x	s3x	s3x	s3x	s3x	s3x	s3x	s3x	s3x
	tramp	vagabundo	бродяга	Tramp	vagabondo	clochard	Tramp	vagabundo	gelandangan	浮浪者	도보 여행가	padyak	आवारा
	condon	perdão	оправдывать	Kondom	condono	pardonner	prezerwatywa	perdonar	memaafkan	コンドン	용납하다	condon	कांडों
	sexs	sexos	пол	Geschlechter	sessi	sexes	płcie	sexos	jenis kelamin	性別	섹스	mga kasarian	लिंग
	adultery	adultério	прелюбодеяние	Ehebruch	adulterio	adultère	cudzołóstwo	adulterio	zina	姦通	간통	pangangalunya	व्यभिचार
	pika	pika	пищуха	Pika	pika	pika	pika	pika	pika	ナキウサギ	피카	pika	पिका
	fck	foda-se	ебать	scheiße	cazzo	putain	kurwa	joder	sialan	ファック	fck	fck	फ़क
	masterbation	masterbação	мастербация	Masterbation	masterbation	masterbation	masterbacja	masterbación	masterbasi	マスターベーション	마스터베이션	masterbation	हस्तमैथुन
	pinche	beliscar	щепотка	kneifen	pizzico	pincer	uszczypnąć	pinchar	mencubit	ピンチェ	핀치	kurot	चुटकी
	bicho	bicho	бичо	bicho	bicho	bicho	bicho	bicho	bicho	ビチョ	비초	bicho	बिचो
	bich	bichinho	сука	Bich	cagna	biche	bicz	bicho	jalang	ビッチ	비치	bich	बिच
	booby	bobo	олух	Dummkopf	tonto	nigaud	głupek	bobo	bodoh	ブービー	얼간이	booby	मूर्ख
	asss	bunda	задница	Arsch	culo	cul	tyłek	culo	pantat	お尻	엉덩이	asss	गधा
	chum	amigo	приятель	Kumpel	amico	ami	kolega	amigo	sahabat karib	友達	친구	chum	दोस्त
	lmfao	lmfao	Имфао	lmfao	lmfao	lmfao	nieważne	lmfao	lmfao	ふむふむ	ㅋㅋㅋㅋㅋ	lmfao	lmfao
	dik	idiota	дик	dik	dick	dik	dik	dik	dik	ディク	딕	dik	डिक
	kamasutra	Kamasutra	камасутра	Kamasutra	Kamasutra	kamasutra	kamasutra	kamasutra	kamasutra	カーマスートラ	카마수트라	kamasutra	कामसूत्र
	titi	titi	Тити	titi	titi	titi	titi	tití	titi	ティティ	티티	titi	कॉक
	raper	estuprador	насильник	Vergewaltiger	stupratore	violeur	raper	violador	pemerkosa	レイパー	강간범	raper	रैपर
	fuk	merda	блядь	Scheiße	cavolo	putain	fuj	joder	fuk	クソ	씨발	fuk	फ़क
	motherfuckers	filhos da puta	ублюдки	Wichser	figli di puttana	enfoirés	skurwysyny	hijos de puta	bajingan	クソ野郎	개자식들	mga bastos	कमीने
	vergin	virgem	девственница	Anfang	vergine	marge	vergin	virgen	hampir	ヴァージン	버진	vergin	vergin
	kurwa	Kurwa	курва	Kurwa	Kurwa	kurwa	kurwa	kurwa	kurwa	曲輪	쿠르와	kurwa	कुरवा
	pecker	pica-pau	клюв	Schwanz	uccello	quéquette	kutas	pájaro carpintero	cotok	ペッカー	곡괭이	pecker	चोंच
	pusy	chato	киска	Pusy	pusì	pusy	cipka	pusy	pus	マンコ	보지	puki	बिल्ली
	mamilo	Mamilo	Мамило	Mamilo	mamilo	maman	mamo	mamilo	mamalo	マミロ	마밀로	mamilo	मैमिलो
	bumb	vagabundo	тупица	Blödmann	stupido	idiot	głupek	vago	bodoh	ダメだ	멍청이	bumb	बम्ब
	nippel	mamilo	сосок	nippel	capezzolo	mamelon	sutek	pezón	puting	ニッペル	니펠	utong	निपल
	pennis	pênis	пеннис	pennis	pene	pennis	pens	pene	uang receh	ペニス	페니스	pennis	लिंग`

var suicideKeywordsText = `自杀	Kill myself	Me matar	Убить себя	Mich umbringen	Uccidermi	Me tuer	Zabiję się	Suicidarme	Bunuh diriku	自分を殺してください	자살	Patayin ang sarili	अपने आप को मारना
自杀	Suicide	Suicídio	Самоубийство	Selbstmord	Suicidio	Suicide	Samobójstwo	Suicidio	Bunuh diri	自殺	자살	Pagpapakamatay	आत्मघाती
自残	Self-harm	Auto-mutilação	Причинять себе вред	Selbstverletzung	Autolesionismo	L'automutilation	Samookaleczenia	Autolesiones	Menyakiti diri sendiri	自傷	자해	Pananakit sa sarili	खुद को नुकसान
自我伤害	Self-injury	Auto ferimento	Членовредительство	Selbstverletzung	Ferita autoinflitta	Automutilation	Samookaleczenie	Auto lastimarse	Cedera diri	自己傷害	자해	Pananakit sa sarili	खुद को चोट
有自杀倾向的	Suicidal	Suicida	Суицидальный	Lebensmüde	Suicida	Suicidaire	Samobójczy	Suicida	Kecenderungan bunuh diri	自殺願望のある	자살	Nagpapakamatay	आत्मघात`

var politicsKeywordsText = `纳粹	nazi	nazista	нацист	Nazi	nazista	nazi	nazi	nazi	Nazi	ナチス	나치	nazi	नाजी
纳粹（复数）	nazis	nazistas	нацисты	Nazis	nazisti	nazis	naziści	nazis	Nazi	ナチス	나치	mga nazi	नाजियों
纳粹分子	nazista	nazista	нацист	Nazist	nazista	nazi	nazista	nazista	Nazista	ナジスタ	나치스타	nazista	नाजिस्ता
纳粹主义	nazism	nazismo	нацизм	Nazismus	nazismo	nazisme	nazizm	nazismo	nazisme	ナチズム	국가 사회주의	nazismo	फ़ासिज़्म
希特勒	hitler	Hitler	Гитлер	Hitler	Hitler	Hitler	hitler	hitler	hitler	ヒトラー	히틀러	hitler	हिटलर
法西斯主义	fascism	fascismo	фашизм	Faschismus	fascismo	fascisme	faszyzm	fascismo	fasisme	ファシズム	파시즘	pasismo	फ़ैसिस्टवाद
法西斯主义（德语）	faschismus	fascismo	фашизм	Faschismus	fascismo	fascisme	faszizm	fascismo	fakismus	ファシスムス	파시즘	faschismus	faschismus
以色列	Israel	Israel	Израиль	Israel	Israele	Israël	Izrael	Israel	Israel	イスラエル	이스라엘	Israel	इजराइल
哈马斯	Hamas	Hamas	ХАМАС	Hamas	Hamas	Hamas	Hamas	Hamás	Hamas	ハマス	하마스	Hamas	हमास
巴勒斯坦	Palestine	Palestina	Палестина	Palästina	Palestina	Palestine	Palestyna	Palestina	Palestina	パレスチナ	팔레스타인	Palestine	फिलिस्तीन
亚伦·布什内尔	Aaron Bushnell	Aaron Bushnell	Аарон Бушнелл	Aaron Bushnell	Aaron Bushnell	Aaron Bushnell	Aarona Bushnella	Aaron Bushnell	Aaron Bushnell	アーロン・ブッシュネル	아론 부시넬	Aaron Bushnell	एरोन बुशनेल`

var underageKeywordsText = `	minors	menores	несовершеннолетние	Minderjährige	minori	mineurs	nieletni	menores	anak di bawah umur	未成年	미성년자	mga menor de edad	नाबालिगों
	underage	menor de idade	несовершеннолетний	minderjährig	minorenne	mineur	niepełnoletni	menor de edad	di bawah umur	未成年	부족	menor de edad	अवयस्क
	1 year old	1 ano	1 год	1 Jahr alt	1 anno	1 an	1 rok	1 año	1 tahun	1歳	1세	1 taong gulang	1 साल का  
	2 year old	2 anos	2 года	2 Jahre alt	2 anni	2 ans	2 lata	2 años	2 tahun	2歳	2세	2 taong gulang	2 साल का  
	3 year old	3 anos	3 года	3 Jahre alt	3 anni	3 ans	3 lata	3 años	3 tahun	3歳	3세	3 taong gulang	3 साल का  
	4 year old	4 anos	4 года	4 Jahre alt	4 anni	4 ans	4 lata	4 años	4 tahun	4歳	4세	4 na taong gulang	4 साल का  
	5 year old	5 anos	5 лет	5 Jahre alt	5 anni	5 ans	5 lat	5 años	5 tahun	5歳	5세	5 taong gulang	5 साल का  
	6 year old	6 anos	6 лет	6 Jahre alt	6 anni	6 ans	6 lat	6 años	6 tahun	6歳	6세	6 taong gulang	6 साल का  
	7 year old	7 anos	7 лет	7 Jahre alt	7 anni	7 ans	7 lat	7 años	7 tahun	7歳	7세	7 taong gulang	7 साल का  
	8 year old	8 anos	8 лет	8 Jahre alt	8 anni	8 ans	8 lat	8 años	8 tahun	8歳	8세	8 taong gulang	8 साल का  
	9 year old	9 anos	9 лет	9 Jahre alt	9 anni	9 ans	9 lat	9 años	9 tahun	9歳	9세	9 taong gulang	9 साल का  
	10 year old	10 anos	10 лет	10 Jahre alt	10 anni	10 ans	10 lat	10 años	10 tahun	10歳	10세	10 taong gulang	10 साल का  
	11 year old	11 anos	11 лет	11 Jahre alt	11 anni	11 ans	11 lat	11 años	11 tahun	11歳	11세	11 taong gulang	11 साल का  
	12 year old	12 anos	12 лет	12 Jahre alt	12 anni	12 ans	12 lat	12 años	12 tahun	12歳	12세	12 taong gulang	12 साल का  
	13 year old	13 anos	13 лет	13 Jahre alt	13 anni	13 ans	13 lat	13 años	13 tahun	13歳	13세	13 taong gulang	13 साल का  
	14 year old	14 anos	14 лет	14 Jahre alt	14 anni	14 ans	14 lat	14 años	14 tahun	14歳	14세	14 taong gulang	14 साल का  
	15 year old	15 anos	15 лет	15-jährige	15 anni	15 ans	15 lat	15 años	15 tahun	15歳	15세	15 taong gulang	15 साल का  
	16 year old	16 anos	16 лет	16 Jahre alt	16 anni	16 ans	16 lat	16 años	16 tahun	16歳	16세	16 taong gulang	16 साल का  
	17 year old	17 anos	17 лет	17 Jahre alt	17 anni	17 ans	17 lat	17 años	17 tahun	17歳	17세	17 taong gulang	17 साल का  
	1 years old	1 ano	1 год	1 Jahr alt	1 anno	1 ans	1 rok	1 años	1 tahun	1歳	1세	1 taong gulang	1 साल का  
	2 years old	2 anos	2 года	2 Jahre alt	2 anni	2 ans	2 lata	2 años	2 tahun	2歳	2세	2 taong gulang	2 साल का  
	3 years old	3 anos	3 года	3 Jahre alt	di 3 anni	3 ans	3 lata	3 años	3 tahun	3歳	3세	3 taong gulang	3 साल का  
	4 years old	4 anos	4 года	4 Jahre alt	4 anni	4 ans	4 lata	4 años	4 tahun	4歳	4세	4 na taong gulang	4 साल का  
	5 years old	5 anos	5 лет	5 Jahre alt	5 anni	5 ans	5 lat	5 años	5 tahun	5歳	5 세	5 taong gulang	5 साल का  
	6 years old	6 anos	6 лет	6 Jahre alt	6 anni	6 ans	6 lat	6 años	6 tahun	6歳	6세	6 taong gulang	6 साल का  
	7 years old	7 anos	7 лет	7 Jahre alt	7 anni	7 ans	7 lat	7 años	7 tahun	7歳	7세	7 taong gulang	7 साल का  
	8 years old	8 anos	8 лет	8 Jahre alt	8 anni	8 ans	8 lat	8 años	8 tahun	8歳	8 살	8 taong gulang	8 साल का  
	9 years old	9 anos	9 лет	9 Jahre alt	9 anni	9 ans	9 lat	9 años	9 tahun	9歳	9세	9 taong gulang	9 साल का  
	10 years old	10 anos	10 лет	10 Jahre alt	10 anni	10 ans	10 lat	10 años	10 tahun	10歳	10세	10 taong gulang	10 साल का  
	11 years old	11 anos	11 лет	11 Jahre alt	11 anni	11 ans	11 lat	11 años	11 tahun	11歳	11 살	11 taong gulang	11 साल का  
	12 years old	12 anos	12 лет	12 Jahre alt	12 anni	12 ans	12 lat	12 años	12 tahun	12歳	12 살	12 taong gulang	12 साल का  
	13 years old	13 anos	13 лет	13 Jahre alt	13 anni	13 ans	13 lat	13 años	13 tahun	13歳	13 살	13 taong gulang	13 साल का  
	14 years old	14 anos	14 лет	14 Jahre alt	14 anni	14 ans	14 lat	14 años	14 tahun	14歳	14 살	14 taong gulang	14 साल का  
	15 years old	15 anos	15 лет	15 Jahre alt	15 anni	15 ans	15 lat	15 años	15 tahun	15歳	15세	15 taong gulang	15 साल का  
	16 years old	16 anos	16 лет	16 Jahre alt	16 anni	16 ans	16 lat	16 años	16 tahun	16歳	16세	16 taong gulang	16 साल का  
	17 years old	17 anos	17 лет	17 Jahre alt	17 anni	17 ans	17 lat	17 años	17 tahun	17歳	17 살	17 taong gulang	17 साल का  
	age1	idade1	возраст1	Alter1	età1	âge1	wiek1	edad1	umur1	1歳	나이1	edad1	उम्र1
	age2	idade2	возраст2	Alter2	età2	âge2	wiek2	edad2	umur2	2歳	나이2	edad2	आयु2
	age3	3 anos	возраст3	Alter3	età3	âge3	wiek3	edad3	usia3	3歳	나이3	edad3	आयु3
	age4	4 anos	возраст4	Alter4	età4	âge4	wiek4	edad4	usia4	4歳	나이4	edad4	उम्र4
	age5	5 anos	возраст5	Alter5	età5	âge5	wiek5	edad5	usia5	5歳	나이5	edad5	आयु5
	age6	6 anos	возраст6	Alter6	età6	6 ans	wiek6	edad6	usia6	6歳	나이6	edad6	उम्र6
	age7	7 anos	возраст7	Alter7	età7	7 ans	wiek7	edad7	usia7	7歳	나이7	edad7	उम्र7
	age8	8 anos	возраст8	Alter8	età8	8 ans	wiek8	edad8	usia8	8歳	나이8	edad8	आयु8
	age9	9 anos	возраст9	Alter9	età9	9 ans	wiek9	edad9	usia9	9歳	나이9	edad9	आयु9
	age10	10 anos	возраст10	Alter10	età10	10 ans	wiek10	edad10	usia10	10歳	나이10	edad10	उम्र10
	age11	11 anos	возраст11	Alter11	età11	11 ans	wiek11	edad11	usia11	11歳	나이11	edad11	उम्र11
	age12	12 anos	возраст12	Alter12	età12	12 ans	wiek12	edad12	usia12	12歳	나이12	edad12	उम्र12
	age13	13 anos	возраст13	Alter13	età13	13 ans	wiek13	edad13	usia13	13歳	나이13	edad13	उम्र13
	age14	14 anos	возраст14	Alter14	età14	14 ans	wiek 14	edad14	usia14	14歳	나이14	edad14	उम्र14
	age15	15 anos	возраст15	Alter15	età15	15 ans	wiek15	edad15	usia15	15歳	나이15	edad15	उम्र15
	age16	16 anos	возраст16	Alter16	età16	16 ans	wiek 16	edad16	usia16	16歳	나이16	edad16	उम्र16
	age17	17 anos	возраст17	Alter17	età17	17 ans	wiek17	edad17	usia17	17歳	나이17	edad17	उम्र17
	age 1	1 ano	возраст 1	Alter 1	età 1	1 an	wiek 1	edad de 1 año	usia 1	1歳	1세	edad 1	उम्र 1
	age 2	2 anos	возраст 2	Alter 2	età 2	2 ans	wiek 2	edad de 2 años	usia 2	2歳	2세	edad 2	उम्र 2
	age 3	3 anos	возраст 3	Alter 3	età 3	3 ans	wiek 3	edad de 3 años	usia 3	3歳	3세	edad 3	उम्र 3
	age 4	4 anos	возраст 4	Alter 4	età 4	4 ans	wiek 4	edad de 4 años	usia 4	4歳	4세	edad 4	उम्र 4
	age 5	5 anos	возраст 5	Alter 5	età 5	5 ans	wiek 5	edad de 5 años	usia 5	5歳	5세	edad 5	उम्र 5
	age 6	6 anos	возраст 6	Alter 6	età 6	6 ans	wiek 6	edad de 6 años	usia 6	6歳	6세	edad 6	उम्र 6
	age 7	7 anos	возраст 7	Alter 7	età 7	7 ans	wiek 7	edad de 7 años	usia 7	7歳	7세	edad 7	उम्र 7
	age 8	8 anos	возраст 8	Alter 8	età 8	8 ans	wiek 8	edad de 8 años	usia 8	8歳	8세	edad 8	उम्र 8
	age 9	9 anos	возраст 9	Alter 9	età 9	9 ans	wiek 9	edad de 9 años	usia 9	9歳	9세	edad 9	उम्र 9
	age 10	10 anos	возраст 10	Alter 10	età 10	10 ans	wiek 10	edad de 10 años	usia 10	10歳	10세	edad 10	उम्र 10
	age 11	11 anos	возраст 11	Alter 11	età  11	11 ans	wiek 11	edad de 11 años	usia 11	11歳	11세	edad 11	उम्र 11
	age 12	12 anos	возраст 12	Alter 12	età 12	12 ans	wiek 12	edad de 12 años	usia 12	12歳	12세	edad 12	उम्र 12
	age 13	13 anos	возраст 13	Alter 13	età 13	13 ans	wiek 13	edad de 13 años	usia 13	13歳	13세	edad 13	उम्र 13
	age 14	14 anos	возраст 14	Alter 14	età 14	14 ans	wiek 14	edad de 14 años	usia 14	14歳	14세	edad 14	उम्र 14
	age 15	15 anos	возраст 15	Alter 15	età 15	15 ans	wiek 15	edad de 15 años	usia 15	15歳	15세	edad 15	उम्र 15
	age 16	16 anos	возраст 16	Alter 16	età 16	16 ans	wiek 16	edad de 16 años	usia 16	16歳	16세	edad 16	उम्र 16
	age 17	17 anos	возраст 17	Alter 17	età 17	17 ans	wiek 17	edad de 17 años	usia 17	17歳	17세	edad 17	उम्र 17
	age=1	idade=1	возраст=1	Alter=1	età=1	âge=1	wiek=1	edad=1	umur=1	年齢=1	나이=1	edad=1	आयु=1
	age=2	idade=2	возраст=2	Alter=2	età=2	âge=2	wiek=2	edad=2	umur=2	年齢=2	나이=2	edad=2	आयु=2
	age=3	idade=3	возраст=3	Alter=3	età=3	âge=3	wiek=3	edad=3	umur=3	年齢=3	나이=3	edad=3	आयु=3
	age=4	idade=4	возраст=4	Alter=4	età=4	âge=4	wiek=4	edad=4	umur=4	年齢=4	나이=4	edad=4	आयु=4
	age=5	idade=5	возраст=5	Alter=5	età=5	âge=5	wiek=5	edad=5	umur=5	年齢=5	나이=5	edad=5	उम्र=5
	age=6	idade=6	возраст=6	Alter=6	età=6	âge=6	wiek=6	edad=6	umur=6	年齢=6	나이=6	edad=6	उम्र=6
	age=7	idade=7	возраст=7	Alter=7	età=7	âge=7	wiek=7	edad=7	umur=7	年齢=7	나이=7	edad=7	आयु=7
	age=8	idade=8	возраст=8	Alter=8	età=8	âge=8	wiek=8	edad=8	umur=8	年齢=8	나이=8	edad=8	उम्र=8
	age=9	idade=9	возраст=9	Alter=9	età=9	âge=9	wiek=9	edad=9	umur=9	年齢=9	나이=9	edad=9	उम्र=9
	age=10	idade=10	возраст=10	Alter=10	età=10	âge=10	wiek=10	edad=10	umur=10	年齢=10	나이=10	edad=10	उम्र=10
	age=11	idade=11	возраст=11	Alter=11	età=11	âge=11	wiek=11	edad=11	umur=11	年齢=11	나이=11	edad=11	उम्र=11
	age=12	idade=12	возраст=12	Alter=12	età=12	âge=12	wiek=12	edad=12	umur=12	年齢=12	나이=12	edad=12	उम्र=12
	age=13	idade=13	возраст=13	Alter=13	età=13	âge=13	wiek=13	edad=13	umur=13	年齢=13	나이=13	edad=13	उम्र=13
	age=14	idade=14	возраст=14	Alter=14	età=14	âge=14	wiek=14	edad=14	umur=14	年齢=14	나이=14	edad=14	उम्र=14
	age=15	idade=15	возраст=15	Alter=15	età=15	âge=15	wiek=15	edad=15	umur = 15	年齢=15	나이=15	edad=15	उम्र=15
	age=16	idade=16	возраст=16	Alter=16	età=16	âge=16	wiek=16	edad=16	umur = 16	年齢=16	나이=16	edad=16	उम्र=16
	age=17	idade=17	возраст=17	Alter=17	età=17	âge=17	wiek=17	edad=17	umur = 17	年齢=17	나이=17	edad=17	उम्र=17
	age of 1	idade de 1	возраст 1 год	Alter von 1	età di 1	l'âge de 1	wiek 1	edad de 1	usia 1 tahun	1歳	1세	edad 1	1 की आयु  
	age of 2	2 anos de idade	возраст 2 года	Alter von 2	età di 2 anni	l'âge de 2 ans	wiek 2 lat	edad de 2	usia 2 tahun	2歳	2세	edad 2	2 की आयु  
	age of 3	idade de 3 anos	возраст 3 года	Alter von 3	età di 3 anni	l'âge de 3 ans	wiek 3 lat	edad de 3	usia 3 tahun	3歳	3세	edad 3	3 की आयु  
	age of 4	4 anos de idade	возраст 4 года	Alter von 4	età di 4 anni	l'âge de 4 ans	wiek 4 lat	edad de 4	usia 4 tahun	4歳	4세	edad 4	4 की आयु  
	age of 5	5 anos de idade	возраст 5 лет	Alter von 5	età di 5 anni	l'âge de 5 ans	wiek 5 lat	edad de 5	usia 5 tahun	5歳	5세	edad 5	5 की आयु  
	age of 6	6 anos de idade	возраст 6 лет	Alter von 6	età di 6 anni	l'âge de 6 ans	wiek 6 lat	edad de 6	usia 6 tahun	6歳	6세	edad 6	6 की आयु  
	age of 7	7 anos de idade	возраст 7 лет	Alter von 7	età di 7 anni	l'âge de 7 ans	wiek 7 lat	edad de 7	usia 7 tahun	7歳	7세	edad 7	7 की आयु  
	age of 8	8 anos de idade	возраст 8 лет	Alter von 8	età di 8 anni	l'âge de 8 ans	wiek 8 lat	edad de 8	usia 8 tahun	8歳	8세	edad 8	8 की आयु  
	age of 9	9 anos de idade	возраст 9 лет	Alter von 9	età di 9 anni	l'âge de 9 ans	wiek 9 lat	edad de 9	usia 9 tahun	9歳	9세	edad 9	9 की आयु  
	age of 10	idade de 10 anos	возраст 10 лет	Alter von 10	età di 10 anni	l'âge de 10 ans	wiek 10 lat	edad de 10	usia 10 tahun	10歳	10세	edad 10	10 की आयु  
	age of 11	11 anos	возраст 11 лет	Alter von 11	età di 11 anni	l'âge de 11 ans	wiek 11 lat	edad de 11	usia 11 tahun	11歳	11세	edad 11	11 की आयु  
	age of 12	12 anos	возраст 12 лет	Alter von 12	età di 12 anni	12 ans	wiek 12 lat	edad de 12	usia 12 tahun	12歳	12세	edad 12	12 की आयु  
	age of 13	13 anos	возраст 13 лет	Alter von 13	età di 13 anni	13 ans	wiek 13 lat	edad de 13	usia 13 tahun	13歳	13세	edad 13	13 की आयु  
	age of 14	14 anos	возраст 14 лет	Alter von 14	età di 14 anni	14 ans	wiek 14 lat	edad de 14	usia 14 tahun	14歳	14세	edad 14	14 की आयु  
	age of 15	15 anos	возраст 15 лет	Alter von 15	età di 15 anni	15 ans	wiek 15 lat	edad de 15	usia 15 tahun	15歳	15세	edad 15	15 की आयु  
	age of 16	16 anos	возраст 16 лет	Alter von 16	età di 16 anni	16 ans	wiek 16 lat	edad de 16	usia 16 tahun	16歳	16세	edad 16	16 की आयु  
	age of 17	17 anos	возраст 17 лет	Alter von 17	età di 17 anni	17 ans	wiek 17 lat	edad de 17	usia 17 tahun	17歳	17세	edad 17	17 की आयु  
	age of one	idade de um	возраст один год	Alter von einem Jahr	età di uno	l'âge d'un an	wiek jednego	edad de un	usia satu tahun	1歳	한 살	age ng isa  	एक वर्ष की आयु
	age of two	dois anos de idade	возраст двух лет	Alter von zwei Jahren	età di due anni	l'âge de deux ans	wiek dwóch	edad de dos	usia dua tahun	2歳	두 살	age ng dalawa  	दो साल की उम्र
	age of three	idade de três	возраст трёх лет	Alter von drei Jahren	età di tre anni	l'âge de trois ans	wiek trzech	edad de tres	usia tiga tahun	3歳	세 살	age ng tatlo  	तीन साल की उम्र
	age of four	quatro anos de idade	возраст четырёх лет	Alter von vier Jahren	età di quattro anni	l'âge de quatre ans	wiek czterech	edad de cuatro	usia empat tahun	4歳	네 살	age ng apat  	चार साल की उम्र
	age of five	cinco anos	возраст пяти лет	Alter von fünf Jahren	età di cinque anni	l'âge de cinq ans	wiek pięciu	edad de cinco	usia lima tahun	5歳	다섯 살	age ng lima  	पांच साल की उम्र
	age of six	seis anos	возраст шести лет	Alter von sechs Jahren	età di sei anni	l'âge de six ans	wiek sześciu	edad de seis	usia enam tahun	6歳	여섯 살	age ng anim  	छह साल की उम्र
	age of seven	idade de sete anos	возраст семи лет	Alter von sieben Jahren	età di sette anni	l'âge de sept ans	wiek siedmiu	edad de siete	usia tujuh tahun	7歳	일곱 살	age ng pito  	सात साल की उम्र
	age of eight	idade de oito anos	возраст восьми лет	Alter von acht Jahren	età di otto anni	l'âge de huit ans	wiek ośmiu	edad de ocho	usia delapan tahun	8歳	여덟 살	age ng walo  	आठ साल की उम्र
	age of nine	nove anos	возраст девяти лет	Alter von neun Jahren	età di nove anni	l'âge de neuf ans	wiek dziewięciu	edad de nueve	usia sembilan tahun	9歳	아홉 살	age ng siyam  	नौ साल की उम्र
	age of ten	dez anos	возраст десяти лет	Alter von zehn Jahren	età di dieci anni	l'âge de dix ans	wiek dziesięciu	edad de diez	usia sepuluh tahun	10歳	열살	age ng sampu  	दस साल की उम्र
	age of eleven	onze anos	возраст одиннадцати лет	Alter von elf Jahren	età di undici anni	l'âge de onze ans	wiek jedenastu	edad de once	usia sebelas tahun	11歳	열한 살	age ng labing-isa  	ग्यारह साल की उम्र
	age of twelve	idade de doze anos	возраст двенадцати лет	Alter von zwölf Jahren	età di dodici anni	l'âge de douze ans	wiek dwunastu	edad de doce	usia dua belas tahun	12歳	열두 살	age ng labing-dalawa  	बारह वर्ष की आयु
	age of thirteen	treze anos	возраст тринадцати лет	Alter von dreizehn Jahren	età di tredici anni	l'âge de treize ans	wiek trzynastu	edad de trece	usia tiga belas tahun	13歳	열세 살	age ng labing-tatlo  	तेरह साल की उम्र
	age of fourteen	quatorze anos	возраст четырнадцати лет	Alter von vierzehn Jahren	età di quattordici anni	l'âge de quatorze ans	wiek czternastu	edad de catorce	usia empat belas tahun	14歳	열네 살	age ng labing-apat  	चौदह वर्ष की आयु
	age of fifteen	quinze anos	возраст пятнадцати лет	Alter von fünfzehn Jahren	età di quindici anni	l'âge de quinze ans	wiek piętnastu	edad de quince	usia lima belas tahun	15歳	열다섯 살	age ng labing-lima  	पंद्रह वर्ष की आयु
	age of sixteen	idade de dezesseis anos	возраст шестнадцати лет	Alter von sechzehn Jahren	età di sedici anni	l'âge de seize ans	wiek szesnastu	edad de dieciséis	usia enam belas tahun	16歳	열여섯 살	age ng labing-anim  	सोलह वर्ष की आयु
	age of seventeen	idade de dezessete anos	возраст семнадцати лет	Alter von siebzehn Jahren	età di diciassette anni	l'âge de dix-sept ans	wiek siedemnastu	edad de diecisiete	usia tujuh belas tahun	17歳	열일곱 살	age ng labing-pito  	सत्रह वर्ष की आयु
	age=one	idade = um	возраст = один	Alter=eins	età=uno	âge = un	wiek = jeden	edad = uno	umur = satu	年齢=1歳	나이=한 살  	edad=isa	उम्र=एक
	age=two	idade = dois	возраст = два	Alter = zwei	età=due	âge = deux	wiek = dwa	edad = dos	umur = dua	年齢=2歳	나이=두 살  	edad=dalawa	उम्र=दो
	age=three	idade = três	возраст = три	Alter = drei	età=tre	âge = trois	wiek = trzy	edad = tres	umur = tiga	年齢=3歳	나이=세 살  	edad=tatlo	उम्र=तीन
	age=four	idade = quatro	возраст=четыре	Alter = vier	età=quattro	âge = quatre	wiek = cztery	edad = cuatro	umur = empat	年齢=4歳	나이=네 살  	edad=apat	उम्र=चार
	age=five	idade = cinco	возраст = пять	Alter = fünf	età=cinque	âge = cinq	wiek = pięć	edad = cinco	umur = lima	年齢=5歳	나이=다섯 살  	edad=lima	उम्र=पांच
	age=six	idade = seis	возраст = шесть	Alter=sechs	età=sei	âge = six	wiek=sześć	edad = seis	umur = enam	年齢=6歳	나이=여섯 살  	edad=anim	उम्र=छह
	age=seven	idade = sete	возраст = семь	Alter = sieben	età=sette	âge = sept	wiek = siedem	edad = siete	umur = tujuh	年齢=7歳	나이=일곱 살  	edad=pito	उम्र=सात
	age=eight	idade = oito	возраст=восемь	Alter = acht	età=otto	âge = huit	wiek = osiem	edad = ocho	umur = delapan	年齢=8歳	나이=여덟 살  	edad=walo	उम्र=आठ
	age=nine	idade = nove	возраст = девять	Alter = neun	età=nove	âge = neuf	wiek = dziewięć	edad = nueve	umur = sembilan	年齢=9歳	나이=아홉 살  	edad=siyam	उम्र=नौ
	age=ten	idade = dez	возраст = десять	Alter=zehn	età=dieci	âge = dix	wiek = dziesięć	edad = diez	umur = sepuluh	年齢=10歳	나이=열 살  	edad=sampu	उम्र=दस
	age=eleven	idade = onze	возраст = одиннадцать	Alter=elf	età=undici	âge = onze	wiek = jedenaście	edad = once	umur=sebelas	年齢=11歳	나이=열한 살  	edad=labingisa	उम्र=ग्यारह
	age=twelve	idade = doze	возраст = двенадцать	Alter=zwölf	età=dodici	âge = douze	wiek = dwanaście	edad = doce	umur = dua	年齢=12歳	나이=열두 살  	edad=labindalawa	उम्र=बारह
	age=thirteen	idade = treze	возраст = тринадцать	Alter = dreizehn	età=tredici	âge = treize	wiek = trzynaście	edad = trece	umur = tiga	年齢=13歳	나이=열세 살  	edad=labing tatlo	उम्र=तेरह
	age=fourteen	idade = quatorze	возраст = четырнадцать	Alter=vierzehn	età=quattordici	âge = quatorze	wiek = czternaście	edad = catorce	umur = empat	年齢=14歳	나이=열네 살  	edad=labing apat	उम्र=चौदह
	age=fifteen	idade = quinze	возраст = пятнадцать	Alter = fünfzehn	età=quindici	âge=quinze	wiek = piętnaście	edad = quince	umur = lima	年齢=15歳	나이=열다섯 살  	edad=labinlima	उम्र=पंद्रह
	age=sixteen	idade = dezesseis	возраст = шестнадцать	Alter=sechzehn	età=sedici	âge=seize	wiek = szesnaście	edad = dieciséis	umur = enam	年齢=16歳	나이=열여섯 살  	edad=labing-anim	उम्र=सोलह
	age=seventeen	idade = dezessete	возраст = семнадцать	Alter = siebzehn	età=diciassette	âge = dix-sept	wiek = siedemnaście	edad = diecisiete	umur = tujuh	年齢=17歳	나이=열일곱 살  	edad=labing pito	उम्र=सत्रह
	age one	um ano	возраст один	ein Jahr	età uno	un an	wiek jeden	un año de edad	usia satu tahun	1歳	나이 한  	edad isa	उम्र एक
	age two	dois anos	возраст два	zwei Jahre	due anni	deux ans	wiek dwa	dos años	usia dua tahun	2歳	나이 두  	edad dalawa	उम्र दो
	age three	três anos	возраст три	drei Jahre	tre anni	trois ans	wiek trzech	tres años	usia tiga tahun	3歳	나이 세  	edad tatlo	उम्र तीन
	age four	quatro anos	возраст четыре	vier Jahre	quattro anni	quatre ans	wiek czterech	cuatro años	usia empat tahun	4歳	나이 네  	edad apat	उम्र चार
	age five	cinco anos	возраст пять	fünf Jahre	cinque anni	cinq ans	wiek pięciu	cinco años	usia lima tahun	5歳	나이 다섯  	edad lima	उम्र पांच
	age six	seis anos	возраст шесть	sechs Jahre	sei anni	six ans	wiek sześciu	seis años	usia enam tahun	6歳	나이 여섯  	edad anim	उम्र छह
	age seven	sete anos	возраст семь	sieben Jahre	sette anni	sept ans	wiek siedmiu	siete años	usia tujuh tahun	7歳	나이 일곱  	edad pito	उम्र सात
	age eight	oito anos	возраст восемь	acht Jahre	otto anni	huit ans	wiek ośmiu	ocho años	usia delapan tahun	8歳	나이 여덟  	edad walo	उम्र आठ
	age nine	nove anos	возраст девять	neun Jahre	nove anni	neuf ans	wiek dziewięciu	nueve años	usia sembilan tahun	9歳	나이 아홉  	edad siyam	उम्र नौ
	age ten	dez anos	возраст десять	zehn Jahre	dieci anni	dix ans	wiek dziesięciu	diez años	usia sepuluh tahun	10歳	나이 열  	edad sampu	उम्र दस
	age eleven	onze anos	возраст одиннадцать	elf Jahre	undici anni	onze ans	jedenaście	once años	umur sebelas	11歳	나이 열한  	labing-isang taong gulang	उम्र ग्यारह
	age twelve	doze anos	возраст двенадцать	zwölf Jahre	dodici anni	douze ans	wiek dwunastu	doce años	usia dua belas tahun	12歳	나이 열두  	edad labindalawa	उम्र बारह
	age thirteen	treze anos	возраст тринадцать	dreizehn Jahre	tredici anni	treize ans	wiek trzynastu	trece años	usia tiga belas tahun	13歳	나이 열세  	edad labintatlo	उम्र तेरह
	age fourteen	quatorze anos	возраст четырнадцать	vierzehn Jahre	quattordici anni	quatorze ans	wiek czternastu	catorce años	usia empat belas tahun	14歳	나이 열네  	edad labing-apat	उम्र चौदह
	age fifteen	quinze anos	возраст пятнадцать	fünfzehn Jahre	quindici anni	quinze ans	wiek piętnastu	quince años	usia lima belas tahun	15歳	나이 열다섯  	edad labinlima	उम्र पंद्रह
	age sixteen	dezesseis anos	возраст шестнадцать	sechzehn Jahre	sedici anni	seize ans	wiek szesnastu	dieciséis años	usia enam belas tahun	16歳	나이 열여섯  	edad labing-anim	उम्र सोलह
	age seventeen	dezessete anos	возраст семнадцать	siebzehn Jahre	diciassette anni	dix-sept ans	wiek siedemnastu	diecisiete años	usia tujuh belas tahun	17歳	나이 열일곱	edad labing pito	उम्र सत्रह
	one year old	um ano de idade	однолетний	ein Jahr alt	un anno	Un an	roczny	un año de edad	Umur satu tahun	1歳	한 살	isang taong gulang	एक साल का  
	two year old	dois anos de idade	двухлетний	zwei Jahre alt	due anni	deux ans	dwulatek	dos años de edad	berumur dua tahun	2歳	두 살	dalawang taong gulang	दो साल का  
	three year old	três anos de idade	трёхлетний	drei Jahre alt	tre anni	trois ans	trzylatek	tres años	berumur tiga tahun	3歳	세 살	tatlong taong gulang	तीन साल का  
	four year old	quatro anos de idade	четырёхлетний	vier Jahre alt	quattro anni	quatre ans	czterolatek	cuatro años	berumur empat tahun	4歳	네 살	apat na taong gulang	चार साल का  
	five year old	cinco anos de idade	пятилетний	fünf Jahre alt	cinque anni	cinq ans	pięciolatek	cinco años de edad	berumur lima tahun	5歳	다섯 살	limang taong gulang	पांच साल का  
	six year old	seis anos de idade	шестилетний	sechs Jahre alt	sei anni	six ans	sześciolatek	de seis años	berusia enam tahun	6歳	여섯 살	anim na taong gulang	छः साल का  
	seven year old	sete anos de idade	семилетний	sieben Jahre alt	sette anni	sept ans	siedmiolatek	siete años	berumur tujuh tahun	7歳	일곱 살	pitong taong gulang	सात साल का  
	eight year old	oito anos de idade	восьмилетний	acht Jahre alt	otto anni	huit ans	ośmiolatek	ocho años	berusia delapan tahun	8歳	여덟 살	walong taong gulang	आठ साल का  
	nine year old	nove anos de idade	девятилетний	neun Jahre alt	nove anni	neuf ans	dziewięć	nueve años	berusia sembilan tahun	9歳	아홉 살	siyam na taong gulang	नौ साल का  
	ten year old	dez anos de idade	десятилетний	zehn Jahre alt	dieci anni	dix ans	dziesięcioletni	diez años	berumur sepuluh tahun	10歳	열 살	sampung taong gulang	दस साल का  
	eleven year old	onze anos de idade	одиннадцатилетний	elf Jahre alt	undici anni	onze ans	jedenastolatek	once años	berumur sebelas tahun	11歳	열한 살	labing-isang taong gulang	ग्यारह साल का  
	twelve year old	doze anos de idade	двенадцатилетний	zwölf Jahre alt	dodici anni	douze ans	dwunastolatek	doce años	berumur dua belas tahun	12歳	열두 살	labindalawang taong gulang	बारह साल का  
	thirteen year old	treze anos de idade	тринадцатилетний	dreizehn Jahre alt	tredici anni	treize ans	trzynastoletni	trece años de edad	berusia tiga belas tahun	13歳	열세 살	labing tatlong taong gulang	तेरह साल का  
	fourteen year old	quatorze anos de idade	четырнадцатилетний	vierzehn Jahre alt	quattordici anni	quatorze ans	czternastolatek	catorce años	berumur empat belas tahun	14歳	열네 살	labing-apat na taong gulang	चौदह साल का  
	fifteen year old	quinze anos de idade	пятнадцатилетний	fünfzehn Jahre alt	quindici anni	quinze ans	piętnastoletni	quince años	berumur lima belas tahun	15歳	열다섯 살	labinlimang taong gulang	पंद्रह साल का  
	sixteen year old	dezesseis anos de idade	шестнадцатилетний	sechzehn Jahre alt	sedici anni	seize ans	szesnastolatek	dieciséis años	enam belas tahun	16歳	열여섯 살	labing-anim na taong gulang	सोलह साल का  
	seventeen year old	dezessete anos de idade	семнадцатилетний	siebzehn Jahre alt	diciassette anni	dix-sept ans	siedemnastoletni	diecisiete años	berumur tujuh belas tahun	17歳	열일곱 살	labing pitong taong gulang	सत्रह साल का
	one years old	um ano de idade	один год	ein Jahr alt  	un anno	un an	jednoroczny  	Un año de edad	satu tahun	1歳	한 살	isang taong gulang	एक साल की उम्र
	two years old	dois anos de idade	два года	zwei Jahre alt  	due anni	deux ans	dwuletni  	dos años	dua tahun	2歳	두 살  	dalawang taong gulang	दो वर्षीय
	three years old	três anos de idade	три года	drei Jahre alt  	tre anni	trois ans	trzyletni  	tres años	tiga tahun	3歳	세 살  	tatlong taong gulang	तीन साल पुराना
	four years old	Quatro anos de idade	четыре года	vier Jahre alt  	di quattro anni	quatre ans	czterolatek  	cuatro años	empat tahun	4歳	네 살  	Apat na taong gulang	चार वर्ष पुराना
	five years old	cinco anos de idade	пять лет	fünf Jahre alt  	cinque anni	cinq ans	pięcioletni  	cinco años	lima tahun	5歳	다섯 살  	limang taong gulang	पांच वर्षीय
	six years old	Seis anos de idade	шесть лет	sechs Jahre alt  	sei anni	six ans	sześcioletni  	seis años	enam tahun	6歳	여섯 살  	anim na taong gulang	छः वर्ष का
	seven years old	sete anos de idade	семь лет	sieben Jahre alt  	sette anni	sept ans	siedmioletni  	siete años	tujuh tahun	7歳	일곱 살  	pitong taong gulang	सात साल की उम्र
	eight years old	oito anos de idade	восемь лет	acht Jahre alt  	otto anni	huit ans	ośmioletni  	ocho años	delapan tahun	8歳	여덟 살  	walong taong gulang	आठ साल का
	nine years old	nove anos de idade	девять лет	neun Jahre alt  	di nove anni	neuf ans	dziewięcioletni  	nueve años	sembilan tahun	9歳	아홉 살  	siyam na taong gulang	नौ साल की उम्र
	ten years old	dez anos de idade	десять лет	zehn Jahre alt  	di dieci anni	dix ans	dziesięcioletni  	diez años	sepuluh tahun	10歳	열 살  	sampung taong gulang	दस साल पुराना
	eleven years old	onze anos de idade	одиннадцать лет	elf Jahre alt  	undici anni	onze ans	jedenastoletni  	once años	sebelas tahun	11歳	열한 살  	labing-isang taong gulang	ग्यारह वर्ष की उम्र
	twelve years old	doze anos de idade	двенадцать лет	zwölf Jahre alt  	dodici anni	douze ans	dwunastoletni  	doce años	dua belas tahun	12歳	열두 살  	labindalawang taong gulang	बारह साल की उम्र
	thirteen years old	Treze anos de idade	тринадцать лет	dreizehn Jahre alt  	tredici anni	treize ans	trzynastoletni  	trece años	tiga belas tahun	13歳	열세 살  	labing tatlong taong gulang	तेरह साल की उम्र
	fourteen years old	quatorze anos de idade	четырнадцать лет	vierzehn Jahre alt  	quattordici anni	quatorze ans	czternastoletni  	catorce años	empat belas tahun	14歳	열네 살  	labing apat na taong gulang	चौदह साल का किशोर
	fifteen years old	quinze anos de idade	пятнадцать лет	fünfzehn Jahre alt  	quindici anni	quinze ans	piętnastoletni  	quince años	lima belas tahun	15歳	열다섯 살  	labinlimang taong gulang	पंद्रह साल की उम्र
	sixteen years old	dezesseis anos de idade	шестнадцать лет	sechzehn Jahre alt	sedici anni	seize ans	szesnastoletni  	dieciséis años	enam belas tahun	16歳	열여섯 살  	labing anim na taong gulang	सोलह साल की आयु
	seventeen years old	Dezessete anos de idade	семнадцать лет	siebzehn Jahre alt	diciassette anni	dix-sept ans	siedemnastoletni  	diecisiete años	tujuh belas tahun	17歳	열일곱 살	labing pitong taong gulang	सत्रह साल की उम्र
	age="1"	idade = "1"	возраст="1"	Alter="1"	età="1"	âge="1"	wiek="1"	edad="1"	usia = "1"	年齢 = 1"	나이="1"	edad="1"	उम्र='1'
	age="2"	idade = "2"	возраст="2"	Alter="2"	età="2"	âge="2"	wiek="2"	edad="2"	usia = "2"	年齢 = 2	나이="2"	edad="2"	उम्र='2'
	age="3"	idade = "3"	возраст="3"	Alter="3"	età="3"	âge="3"	wiek="3"	edad="3"	usia = "3"	年齢="3"	나이="3"	edad="3"	उम्र='3'
	age="4"	idade = "4"	возраст="4"	Alter="4"	età="4"	âge="4"	wiek="4"	edad="4"	usia = "4"	年齢="4"	나이="4"	edad="4"	उम्र='4'
	age="5"	idade = "5"	возраст="5"	Alter="5"	età="5"	âge="5"	wiek="5"	edad="5"	usia = "5"	年齢="5"	나이="5"	edad="5"	उम्र='5'
	age="6"	idade = "6"	возраст="6"	Alter="6"	età="6"	âge="6"	wiek="6"	edad="6"	usia = "6"	年齢="6"	나이="6"	edad="6"	उम्र='6'
	age="7"	idade = "7"	возраст="7"	Alter="7"	età="7"	âge="7"	wiek="7"	edad="7"	usia = "7"	年齢="7"	나이="7"	edad="7"	उम्र='7'
	age="8"	idade = "8"	возраст="8"	Alter="8"	età="8"	âge="8"	wiek="8"	edad="8"	usia = "8"	年齢="8"	나이="8"	edad="8"	उम्र='8'
	age="9"	idade = "9"	возраст="9"	Alter="9"	età="9"	âge="9"	wiek="9"	edad="9"	usia = "9"	年齢="9"	나이="9"	edad="9"	उम्र='9'
	age="10"	idade = "10"	возраст="10"	Alter="10"	età="10"	âge="10"	wiek="10"	edad="10"	umur = "10"	年齢="10"	나이="10"	edad="10"	उम्र='10'
	age="11"	idade = "11"	возраст="11"	Alter="11"	età="11"	âge="11"	wiek="11"	edad="11"	umur = "11"	年齢 = "11"	나이="11"	edad="11"	उम्र='11'
	age="12"	idade = "12"	возраст="12"	Alter="12"	età="12"	âge="12"	wiek="12"	edad="12"	usia = "12"	年齢 = "12"	나이="12"	edad="12"	उम्र='12'
	age="13"	idade = "13"	возраст="13"	Alter="13"	età="13"	âge="13"	wiek="13"	edad="13"	usia = "13"	年齢 = "13"	나이="13"	edad="13"	उम्र='13'
	age="14"	idade = "14"	возраст="14"	Alter="14"	età="14"	âge="14"	wiek="14"	edad="14"	usia = "14"	年齢 = "14"	나이="14"	edad="14"	उम्र='14'
	age="15"	idade = "15"	возраст="15"	Alter="15"	età="15"	âge="15"	wiek="15"	edad="15"	usia = "15"	年齢 = "15"	나이="15"	edad="15"	उम्र='15'
	age="16"	idade = "16"	возраст="16"	Alter="16"	età="16"	âge="16"	wiek="16"	edad="16"	usia = "16"	年齢 = "16"	나이="16"	edad="16"	उम्र='16'
	age="17"	idade = "17"	возраст="17"	Alter="17"	età="17"	âge="17"	wiek="17"	edad="17"	usia = "17"	年齢 = "17"	나이="17"	edad="17"	उम्र='17'
	age="one"	idade = "um"	возраст="один"	Alter = „eins“	età="uno"	âge = "un"	wiek="jeden"	edad = "uno"	umur = "satu"	年齢=1	나이="한 살"  	edad = "isa"	उम्र='एक'
	age="two"	idade = "dois"	возраст="два"	Alter = „zwei“	età="due"	âge = "deux"	wiek="dwa"	edad = "dos"	umur = "dua"	年齢=2	나이="두 살"  	edad="dalawa"	उम्र='दो'
	age="three"	idade = "três"	возраст="три"	Alter = „drei“	età="tre"	âge = "trois"	wiek="trzy"	edad = "tres"	umur = "tiga"	年齢=3	나이="세 살"  	edad="tatlo"	उम्र='तीन'
	age="four"	idade = "quatro"	возраст="четыре"	Alter = „vier“	età="quattro"	âge = "quatre"	wiek="cztery"	edad = "cuatro"	umur = "empat"	年齢=4	나이="네 살"  	edad="apat"	उम्र='चार'
	age="five"	idade = "cinco"	возраст="пять"	Alter = „fünf“	età="cinque"	âge = "cinq"	wiek="pięć"	edad = "cinco"	umur = "lima"	年齢=5	나이="다섯 살"  	edad="lima"	उम्र='पाँच'
	age="six"	idade = "seis"	возраст="шесть"	Alter = „sechs“	età="sei"	âge = "six"	wiek="sześć"	edad = "seis"	umur = "enam"	年齢=6	나이="여섯 살"  	edad="anim"	उम्र='छह'
	age="seven"	idade = "sete"	возраст="семь"	Alter = „sieben“	età="sette"	âge = "sept"	wiek="siedem"	edad = "siete"	umur = "tujuh"	年齢=7	나이="일곱 살"  	edad="pito"	उम्र='सात'
	age="eight"	idade = "oito"	возраст="восемь"	Alter = „acht“	età="otto"	âge = "huit"	wiek="osiem"	edad = "ocho"	umur = "delapan"	年齢=8	나이="여덟 살"  	edad="eight"	उम्र='आठ'
	age="nine"	idade = "nove"	возраст="девять"	Alter = „neun“	età="nove"	âge = "neuf"	wiek="dziewięć"	edad = "nueve"	umur = "sembilan"	年齢=9	나이="아홉 살"  	edad="siyam"	उम्र='नौ'
	age="ten"	idade = "dez"	возраст="десять"	Alter = „zehn“	età="dieci"	âge = "dix"	wiek="dziesięć"	edad = "diez"	umur = "sepuluh"	年齢=10	나이="열 살"  	edad="sampu"	उम्र=दस'
	age="eleven"	idade = "onze"	возраст="одиннадцать"	Alter = „elf“	età="undici"	âge = "onze"	wiek="jedenaście"	edad = "once"	umur = "sebelas"	年齢=11	나이="열한 살"  	edad="labing isang"	उम्र='ग्यारह'
	age="twelve"	idade = "doze"	возраст="двенадцать"	Alter = „zwölf“	età="dodici"	âge = "douze"	wiek="dwanaście"	edad = "doce"	umur = "dua belas"	年齢=12	나이="열두 살"  	edad="labindalawa"	उम्र='बारह'
	age="thirteen"	idade = "treze"	возраст="тринадцать"	Alter = „dreizehn“	età="tredici"	âge = "treize"	wiek="trzynaście"	edad = "trece"	umur = "tiga belas"	年齢=13	나이="열세 살"  	edad="labintatlo"	उम्र='तेरह'
	age="fourteen"	idade = "quatorze"	возраст="четырнадцать"	Alter = „vierzehn“	età="quattordici"	âge = "quatorze"	wiek="czternaście"	edad = "catorce"	umur = "empat belas"	年齢=14	나이="열네 살"  	edad="labing-apat"	उम्र='चौदह'
	age="fifteen"	idade = "quinze"	возраст="пятнадцать"	Alter = „fünfzehn“	età="quindici"	âge = "quinze"	wiek="piętnaście"	edad = "quince"	umur="lima belas"	年齢=15	나이="열다섯 살"  	edad="labinlima"	उम्र='पंद्रह'
	age="sixteen"	idade = "dezesseis"	возраст="шестнадцать"	Alter = „sechzehn“	età="sedici"	âge = "seize"	wiek="szesnaście"	edad = "dieciséis"	umur = "enam belas"	年齢=16	나이="열여섯 살"  	edad="labing-anim"	उम्र='सोलह'
	age="seventeen"	idade = "dezessete"	возраст="семнадцать"	Alter = „siebzehn“	età="diciassette"	âge = "dix-sept"	wiek="siedemnaście"	edad = "diecisiete"	umur = "tujuh belas"	年齢=17	나이="열일곱 살"  	edad="labing pito"	उम्र='सत्रह'
	age: 1	idade: 1	возраст: 1	Alter: 1	età: 1	âge: 1	wiek: 1	edad: 1	usia: 1	年齢: 1	나이: 1	edad: 1	आयु: 1
	age: 2	idade: 2	возраст: 2	Alter: 2	età: 2	âge: 2	wiek: 2	edad: 2	usia: 2	年齢: 2	나이: 2	edad: 2	आयु: 2
	age: 3	idade: 3	возраст: 3	Alter: 3	età: 3	âge: 3	wiek: 3	edad: 3	usia: 3	年齢: 3	나이: 3	edad: 3	उम्र: 3
	age: 4	idade: 4	возраст: 4	Alter: 4	età: 4	âge: 4	wiek: 4	edad: 4	usia: 4	年齢: 4	나이: 4	edad: 4	उम्र: 4
	age: 5	idade: 5	возраст: 5	Alter: 5	età: 5	âge: 5	wiek: 5	edad: 5	usia: 5	年齢: 5	나이: 5	edad: 5	उम्र: 5
	age: 6	idade: 6	возраст: 6	Alter: 6	età: 6	âge: 6	wiek: 6	edad: 6	usia: 6	年齢: 6	나이: 6	edad: 6	उम्र: 6
	age: 7	idade: 7	возраст: 7	Alter: 7	età: 7	âge: 7	wiek: 7	edad: 7	usia: 7	年齢: 7	나이: 7	edad: 7	उम्र: 7
	age: 8	idade: 8	возраст: 8	Alter: 8	età: 8	âge: 8	wiek: 8	edad: 8	usia: 8	年齢: 8	나이: 8	edad: 8	उम्र: 8
	age: 9	idade: 9	возраст: 9	Alter: 9	età: 9	âge: 9	wiek: 9	edad: 9	usia: 9	年齢: 9	나이: 9	edad: 9	उम्र: 9
	age: 10	idade: 10	возраст: 10	Alter: 10	età: 10	âge: 10	wiek: 10	edad: 10	usia: 10	年齢: 10	나이: 10	edad: 10	उम्र: 10
	age: 11	idade: 11	возраст: 11	Alter: 11	età: 11	âge: 11	wiek: 11	edad: 11	usia: 11	年齢: 11	나이: 11	edad: 11	उम्र: 11
	age: 12	idade: 12	возраст: 12	Alter: 12	età: 12	âge: 12	wiek: 12	edad: 12	usia: 12	年齢: 12	나이: 12	edad: 12	उम्र: 12
	age: 13	idade: 13	возраст: 13	Alter: 13	età: 13	âge: 13	wiek: 13	edad: 13	usia: 13	年齢: 13	나이: 13	edad: 13	उम्र: 13
	age: 14	idade: 14	возраст: 14	Alter: 14	età: 14	âge: 14	wiek: 14	edad: 14	usia: 14	年齢: 14	나이: 14	edad: 14	उम्र: 14
	age: 15	idade: 15	возраст: 15	Alter: 15	età: 15	âge : 15	wiek: 15	edad 15	usia: 15	年齢: 15	나이: 15	edad: 15	उम्र: 15
	age: 16	idade: 16	возраст: 16	Alter: 16	età: 16	âge: 16	wiek: 16	edad: 16	usia: 16	年齢: 16	나이: 16	edad: 16	उम्र: 16
	age: 17	idade: 17	возраст: 17	Alter: 17	età: 17	âge : 17	wiek: 17	edad: 17	usia: 17	年齢: 17	나이: 17	edad: 17	उम्र: 17
	age: one	idade: um	возраст: один	Alter: eins	età: uno	âge: un	wiek: jeden	edad: uno	usia: satu	年齢: 1歳	나이: 하나	edad: isa	उम्र: एक
	age: two	idade: dois	возраст: два	Alter: zwei	età: due	âge: deux	wiek: dwa	edad: dos	usia: dua	年齢: 2歳	나이: 둘	edad: dalawa	उम्र: दो
	age: three	idade: três	возраст: три	Alter: drei	età: tre	âge: trois	wiek: trzy	edad: tres	usia: tiga	年齢: 3歳	나이: 세	edad: tatlo	उम्र: तीन
	age: four	idade: quatro	возраст: четыре	Alter: vier	età: quattro	âge: quatre	wiek: cztery	edad: cuatro	usia: empat	年齢: 4歳	나이: 네	edad: apat	उम्र: चार
	age: five	idade: cinco	возраст: пять	Alter: fünf	età: cinque	âge: cinq	wiek: pięć	edad: cinco	usia: lima	年齢: 5歳	나이: 다섯	edad: lima	उम्र: पांच
	age: six	idade: seis	возраст: шесть	Alter: sechs	età: sei	âge: six	wiek: sześć	edad: seis	usia: enam	年齢: 6歳	나이: 여섯	edad: anim	उम्र: छह
	age: seven	idade: sete	возраст: семь	Alter: sieben	età: sette	âge : sept	wiek: siedem	edad: siete	usia: tujuh	年齢: 7歳	나이: 일곱	edad: pito	उम्र: सात
	age: eight	idade: oito	возраст: восемь	Alter: acht	età: otto	âge: huit	wiek: osiem	edad: ocho	usia: delapan	年齢: 8歳	나이: 여덟	edad: walo	उम्र: आठ
	age: nine	idade: nove	возраст: девять	Alter: neun	età: nove	âge : neuf ans	wiek: dziewięć	edad: nueve	usia: sembilan	年齢: 9歳	나이: 아홉	edad: siyam	उम्र: नौ
	age: ten	idade: dez	возраст: десять	Alter: zehn	età: dieci	âge : dix ans	wiek: dziesięć	edad: diez	usia: sepuluh	年齢: 10歳	나이: 열	edad: sampu	उम्र: दस
	age: eleven	idade: onze	возраст: одиннадцать	Alter: elf	età: undici	âge: onze	wiek: jedenaście	edad: once	usia: sebelas	年齢: 11歳	나이: 열한 살	edad: labing-isa	उम्र: ग्यारह
	age: twelve	idade: doze	возраст: двенадцать	Alter: zwölf	età: dodici	âge : douze	wiek: dwanaście	edad: doce	usia: dua belas	年齢: 12歳	나이: 열두 살	edad: labindalawa	उम्र: बारह
	age: thirteen	idade: treze	возраст: тринадцать	Alter: dreizehn	età: tredici	âge : treize	wiek: trzynaście	edad: trece	usia: tiga belas	年齢: 13歳	나이: 열세 살	edad: labintatlo	उम्र: तेरह
	age: fourteen	idade: quatorze	возраст: четырнадцать	Alter: vierzehn	età: quattordici	âge : quatorze	wiek: czternaście	edad: catorce	usia: empat belas	年齢: 14歳	나이: 열네살	edad: labing-apat	उम्र: चौदह
	age: fifteen	idade: quinze	возраст: пятнадцать	Alter: fünfzehn	età: quindici	âge : quinze	wiek: piętnaście	edad: quince	usia: lima belas	年齢: 15歳	나이: 열다섯	edad: labinlima	उम्र: पंद्रह
	age: sixteen	idade: dezesseis	возраст: шестнадцать	Alter: sechzehn	età: sedici	âge : seize ans	wiek: szesnaście	edad: dieciséis	umur: enam belas	年齢: 16歳	나이: 열여섯	edad: labing-anim	उम्र: सोलह
	age: seventeen	idade: dezessete	возраст: семнадцать	Alter: siebzehn	età: diciassette	âge : dix-sept	wiek: siedemnaście	edad: diecisiete	usia: tujuh belas	年齢: 17歳	나이: 열일곱	edad: labing pito	उम्र: सत्रह
	age"1"	idade"1"	возраст "1"	Alter „1“	età"1"	âge "1"	wiek"1"	edad "1"	umur"1"	年齢「1」	나이"1"	edad "1"	उम्र"1"
	age"2"	idade"2"	возраст "2"	Alter „2“	età"2"	âge "2"	wiek"2"	edad "2"	umur"2"	年齢「2」	나이"2"	edad "2"	उम्र"2"
	age"3"	idade"3"	возраст "3"	Alter „3“	età"3"	âge "3"	wiek"3"	edad "3"	umur"3"	年齢「3」	나이"3"	edad "3"	उम्र"3"
	age"4"	idade"4"	возраст"4"	Alter „4“	età"4"	âge "4"	wiek"4"	edad "4"	umur"4"	年齢「4」	나이"4"	edad "4"	उम्र"4"
	age"5"	idade"5"	возраст"5"	Alter „5“	età"5"	âge "5"	wiek"5"	edad "5"	umur"5"	年齢「5」	나이"5"	edad "5"	उम्र"5"
	age"6"	idade"6"	возраст"6"	Alter"6"	età"6"	âge "6"	wiek"6"	edad "6"	umur"6"	年齢「6」	나이"6"	edad "6"	उम्र"6"
	age"7"	idade"7"	возраст"7"	Alter „7“	età"7"	âge "7"	wiek"7"	edad "7"	umur"7"	年齢「7」	나이"7"	edad "7"	उम्र"7"
	age"8"	idade"8"	возраст"8"	Alter „8“	età"8"	âge "8"	wiek"8"	edad "8"	umur"8"	年齢「8」	나이"8"	edad "8"	उम्र"8"
	age"9"	idade"9"	возраст "9"	Alter „9“	età"9"	âge "9"	wiek"9"	edad "9"	umur"9"	年齢「9」	나이"9"	edad "9"	उम्र"9"
	age"10"	idade"10"	возраст"10"	Alter „10“	età"10"	âge "10"	wiek"10"	edad "10"	umur"10"	年齢「10」	나이"10"	edad "10"	उम्र"10"
	age"11"	idade"11"	возраст"11"	Alter „11“	età"11"	âge "11"	wiek"11"	edad "11"	umur"11"	年齢「11」	나이"11"	edad "11"	उम्र"11"
	age"12"	idade"12"	возраст"12"	Alter „12“	età"12"	âge "12"	wiek"12"	edad "12"	usia "12"	年齢「12」	나이"12"	edad "12"	उम्र"12"
	age"13"	idade"13"	возраст"13"	Alter „13“	età"13"	âge "13"	wiek"13"	edad "13"	umur"13"	年齢「13」	나이"13"	edad "13"	उम्र"13"
	age"14"	idade"14"	возраст"14"	Alter „14“	età"14"	âge "14"	wiek"14"	edad"14"	umur"14"	年齢「14」	나이"14"	edad "14"	उम्र"14"
	age"15"	idade"15"	возраст"15"	Alter „15“	età"15"	âge "15"	wiek"15"	edad 15"	umur"15"	年齢「15」	나이"15"	edad "15"	उम्र"15"
	age"16"	idade"16"	возраст"16"	Alter „16“	età"16"	âge "16"	wiek"16"	edad"16"	umur"16"	年齢「16」	나이"16"	edad "16"	उम्र"16"
	age"17"	idade"17"	возраст"17"	Alter „17“	età"17"	âge "17"	wiek"17"	edad"17"	umur"17"	年齢「17」	나이"17"	edad "17"	उम्र"17"
	age "1"	idade "1"	возраст "1"	Alter „1“	età "1"	âge "1"	wiek „1”	edad "1"	usia "1"	年齢「1」	나이 "1"	edad "1"	उम्र "1"
	age "2"	idade "2"	возраст "2"	Alter „2“	età "2"	âge "2"	wiek „2”	edad "2"	usia "2"	年齢「2」	나이 "2"	edad "2"	उम्र "2"
	age "3"	idade "3"	возраст "3"	Alter „3“	età "3"	âge "3"	wiek „3”	edad "3"	usia "3"	年齢「3」	나이 "3"	edad "3"	उम्र "3"
	age "4"	idade "4"	возраст "4"	Alter „4“	età "4"	âge "4"	wiek „4”	edad "4"	usia "4"	年齢「4」	나이 "4"	edad "4"	उम्र "4"
	age "5"	idade "5"	возраст "5"	Alter „5“	età "5"	âge "5"	wiek „5”	edad "5"	usia "5"	年齢「5」	나이 "5"	edad "5"	उम्र "5"
	age "6"	idade "6"	возраст "6"	Alter „6“	età "6"	âge "6"	wiek „6”	edad "6"	usia "6"	年齢「6」	나이 "6"	edad "6"	उम्र "6"
	age "7"	idade "7"	возраст "7"	Alter „7“	età "7"	âge "7"	wiek „7”	edad "7"	usia "7"	年齢「7」	나이 "7"	edad "7"	उम्र "7"
	age "8"	idade "8"	возраст "8"	Alter „8“	età "8"	âge "8"	wiek „8”	edad "8"	usia "8"	年齢「8」	나이 "8"	edad "8"	उम्र "8"
	age "9"	idade "9"	возраст "9"	Alter „9“	età "9"	âge "9"	wiek „9”	edad "9"	usia "9"	年齢「9」	나이 "9"	edad "9"	उम्र "9"
	age "10"	idade "10"	возраст "10"	Alter „10“	età "10"	âge "10"	wiek „10”	edad "10"	usia "10"	年齢「10」	나이 "10"	edad "10"	उम्र "10"
	age "11"	idade "11"	возраст "11"	Alter „11“	età "11"	âge "11"	wiek „11”	edad "11"	usia "11"	年齢「11」	나이 "11"	edad "11"	उम्र "11"
	age "12"	idade "12"	возраст "12"	Alter „12“	età "12"	âge "12"	wiek „12”	edad "12"	usia "12"	年齢「12」	나이 "12"	edad "12"	उम्र "12"
	age "13"	idade "13"	возраст "13"	Alter „13“	età "13"	âge "13"	wiek „13”	edad "13"	usia "13"	年齢「13」	나이 "13"	edad "13"	उम्र "13"
	age "14"	idade "14"	возраст "14"	Alter „14“	età "14"	âge "14"	wiek „14”	edad "14"	usia "14"	年齢「14」	나이 "14"	edad "14"	उम्र "14"
	age "15"	idade "15"	возраст "15"	Alter „15“	età "15"	âge "15"	wiek „15”	edad 15"	usia "15"	年齢「15」	나이 "15"	edad "15"	उम्र "15"
	age "16"	idade "16"	возраст "16"	Alter „16“	età "16"	âge "16"	wiek „16”	edad "16"	usia "16"	年齢「16」	나이 "16"	edad "16"	उम्र "16"
	age "17"	idade "17"	возраст "17"	Alter „17“	età "17"	âge "17"	wiek „17”	edad "17"	usia "17"	年齢「17」	나이 "17"	edad "17"	उम्र "17"
	age "one"	idade "um"	возраст «один»	Alter „eins“	età "uno"	âge "un"	wiek „jeden”	edad "uno"	usia "satu"	年齢「1」	나이 "하나"	edad "isa"	उम्र "एक"
	age "two"	idade "dois"	возраст «два»	Alter „zwei“	età "due"	âge "deux"	wiek „dwa”	edad "dos"	usia "dua"	年齢「２」	나이 "둘"	edad "dalawa"	उम्र "दो"
	age "three"	idade "três"	возраст «три»	Alter „drei“	età "tre"	âge "trois"	wiek „trzy”	edad "tres"	usia "tiga"	年齢「3」	나이 "세"	edad "tatlo"	उम्र "तीन"
	age "four"	idade "quatro"	возраст «четыре»	Alter „vier“	età "quattro"	âge "quatre"	wiek „cztery”	edad "cuatro"	usia "empat"	年齢「4」	나이 "네"	edad "apat"	उम्र "चार"
	age "five"	idade "cinco"	возраст «пять»	Alter „fünf“	età "cinque"	âge "cinq"	wiek „pięć”	edad "cinco"	usia "lima"	年齢「5」	나이 "다섯"	edad "limang"	उम्र "पांच"
	age "six"	idade "seis"	возраст «шесть»	Alter „sechs“	età "sei"	âge "six"	wiek „sześć”	edad "seis"	usia "enam"	年齢「6」	나이 "여섯"	edad "anim"	उम्र "छह"
	age "seven"	idade "sete"	возраст «семь»	Alter „sieben“	età "sette"	âge "sept"	wiek „siedmiu”	edad "siete"	usia "tujuh"	年齢「7」	나이 "일곱"	edad "pito"	उम्र "सात"
	age "eight"	idade "oito"	возраст «восемь»	Alter „acht“	età "otto"	âge "huit"	wiek „osiem”	edad "ocho"	usia "delapan"	年齢「8」	나이 "여덟"	edad "walong"	उम्र "आठ"
	age "nine"	idade "nove"	возраст «девять»	Alter „neun“	età "nove"	âge "neuf"	wiek „dziewięć”	edad "nueve"	usia "sembilan"	年齢「9」	나이 "아홉"	edad "siyam"	उम्र "नौ"
	age "ten"	idade "dez"	возраст «десять»	Alter „zehn“	età "dieci"	âge "dix"	wiek „dziesięć”	edad "diez"	usia "sepuluh"	年齢「10」	나이 "열"	edad "sampu"	उम्र "दस"
	age "eleven"	idade "onze"	возраст «одиннадцать»	Alter „elf“	età "undici"	âge "onze"	wiek „jedenaście”	edad "once"	usia "sebelas"	年齢「11」	나이 "열한 살"	edad "labing isang"	उम्र "ग्यारह"
	age "twelve"	idade "doze"	возраст «двенадцать»	Alter „zwölf“	età "dodici"	âge "douze"	wiek „dwanaście”	edad "doce"	usia "dua belas"	年齢「12」	나이 "열두 살"	edad "labindalawa"	उम्र "बारह"
	age "thirteen"	idade "treze"	возраст «тринадцать»	Alter „dreizehn“	età "tredici"	âge "treize"	wiek „trzynaście”	edad "trece"	usia "tiga belas"	年齢「13」	나이 "열세 살"	edad "labing tatlo"	उम्र "तेरह"
	age "fourteen"	idade "quatorze"	возраст «четырнадцать»	Alter „vierzehn“	età "quattordici"	âge "quatorze"	wiek „czternaście”	edad "catorce"	usia "empat belas"	年齢「14」	나이 "열네살"	edad "labing-apat"	उम्र "चौदह"
	age "fifteen"	idade "quinze"	возраст «пятнадцать»	Alter „fünfzehn“	età "quindici"	âge "quinze"	wiek „piętnaście”	edad "quince"	usia "lima belas"	年齢「15」	나이 "열다섯"	edad "labinlima"	उम्र "पंद्रह"
	age "sixteen"	idade "dezesseis"	возраст «шестнадцать»	Alter „sechzehn“	età "sedici"	âge "seize"	wiek „szesnaście”	edad "dieciséis"	usia "enam belas"	年齢「16」	나이 "열여섯"	edad "labing-anim"	उम्र "सोलह"
	age "seventeen"	idade "dezessete"	возраст «семнадцать»	Alter „siebzehn“	età "diciassette"	âge "dix-sept"	wiek „siedemnaście”	edad "diecisiete"	usia "tujuh belas"	年齢「17」	나이 "열일곱"	edad "labing pito"	उम्र "सत्रह"`
