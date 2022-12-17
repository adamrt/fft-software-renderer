// This file contains the ability to parse GNS record data.
//
// There is one GNS file per map and an arbitrary number of GNSRecords per map.
// Each record has a type (texture, mesh data, etc). They are always 20 bytes long
// and we read them for a particular map until we hit a RecordTypeEnd.
//
// The GNSRecords contain the location of the data (texture, mesh data, etc) within the
// ISO file. They also contain properties of the map such as time/weather.
package main

import (
	"encoding/binary"
)

type RecordType int

const (
	RecordTypeTexture      RecordType = 0x1701
	RecordTypeMeshPrimary  RecordType = 0x2E01
	RecordTypeMeshOverride RecordType = 0x2F01
	RecordTypeMeshAlt      RecordType = 0x3001
	RecordTypeEnd          RecordType = 0x3101
)

type MapWeather int

const (
	WeatherNone       MapWeather = 0x0
	WeatherNoneAlt    MapWeather = 0x1
	WeatherNormal     MapWeather = 0x2
	WeatherStrong     MapWeather = 0x3
	WeatherVeryStrong MapWeather = 0x4
)

type MapTime int8

const (
	TimeDay   MapTime = 0x0
	TimeNight MapTime = 0x1
)

type GNSRecord []byte

const GNSRecordLen = 20

func (r GNSRecord) Sector() int64 {
	return int64(binary.LittleEndian.Uint16(r[8:10]))
}

func (r GNSRecord) Len() int64 {
	return int64(binary.LittleEndian.Uint32(r[12:16]))
}

func (r GNSRecord) Type() RecordType {
	return RecordType(int(binary.LittleEndian.Uint16(r[4:6])))
}

func (r GNSRecord) Time() MapTime {
	return MapTime(int((r[3] >> 7) & 0x1))
}

func (r GNSRecord) Weather() MapWeather {
	return MapWeather(int((r[3] >> 4) & 0x7))
}

var GNSSectors = [126]int64{
	10026, // MAP000.GNS
	11304, // MAP001.GNS
	12656, // MAP002.GNS
	12938, // MAP003.GNS
	13570, // MAP004.GNS
	14239, // MAP005.GNS
	14751, // MAP006.GNS
	15030, // MAP007.GNS
	15595, // MAP008.GNS
	16262, // MAP009.GNS
	16347, // MAP010.GNS
	16852, // MAP011.GNS
	17343, // MAP012.GNS
	17627, // MAP013.GNS
	18175, // MAP014.GNS
	19510, // MAP015.GNS
	20075, // MAP016.GNS
	20162, // MAP017.GNS
	20745, // MAP018.GNS
	21411, // MAP019.GNS
	21692, // MAP020.GNS
	22270, // MAP021.GNS
	22938, // MAP022.GNS
	23282, // MAP023.GNS
	23557, // MAP024.GNS
	23899, // MAP025.GNS
	23988, // MAP026.GNS
	24266, // MAP027.GNS
	24544, // MAP028.GNS
	24822, // MAP029.GNS
	25099, // MAP030.GNS
	25764, // MAP031.GNS
	26042, // MAP032.GNS
	26229, // MAP033.GNS
	26362, // MAP034.GNS
	27028, // MAP035.GNS
	27643, // MAP036.GNS
	27793, // MAP037.GNS
	28467, // MAP038.GNS
	28555, // MAP039.GNS
	29165, // MAP040.GNS
	29311, // MAP041.GNS
	29653, // MAP042.GNS
	29807, // MAP043.GNS
	30473, // MAP044.GNS
	30622, // MAP045.GNS
	30966, // MAP046.GNS
	31697, // MAP047.GNS
	32365, // MAP048.GNS
	33032, // MAP049.GNS
	33701, // MAP050.GNS
	34349, // MAP051.GNS
	34440, // MAP052.GNS
	34566, // MAP053.GNS
	34647, // MAP054.GNS
	34745, // MAP055.GNS
	35350, // MAP056.GNS
	35436, // MAP057.GNS
	35519, // MAP058.GNS
	35603, // MAP059.GNS
	35683, // MAP060.GNS
	35765, // MAP061.GNS
	36052, // MAP062.GNS
	36394, // MAP063.GNS
	36530, // MAP064.GNS
	36612, // MAP065.GNS
	37214, // MAP066.GNS
	37817, // MAP067.GNS
	38386, // MAP068.GNS
	38473, // MAP069.GNS
	38622, // MAP070.GNS
	39288, // MAP071.GNS
	39826, // MAP072.GNS
	40120, // MAP073.GNS
	40724, // MAP074.GNS
	41391, // MAP075.GNS
	41865, // MAP076.GNS
	42532, // MAP077.GNS
	43200, // MAP078.GNS
	43295, // MAP079.GNS
	43901, // MAP080.GNS
	44569, // MAP081.GNS
	45044, // MAP082.GNS
	45164, // MAP083.GNS
	45829, // MAP084.GNS
	46498, // MAP085.GNS
	47167, // MAP086.GNS
	47260, // MAP087.GNS
	47928, // MAP088.GNS
	48595, // MAP089.GNS
	49260, // MAP090.GNS
	49538, // MAP091.GNS
	50108, // MAP092.GNS
	50387, // MAP093.GNS
	50554, // MAP094.GNS
	51120, // MAP095.GNS
	51416, // MAP096.GNS
	52082, // MAP097.GNS
	52749, // MAP098.GNS
	53414, // MAP099.GNS
	53502, // MAP100.GNS
	53579, // MAP101.GNS
	53659, // MAP102.GNS
	54273, // MAP103.GNS
	54359, // MAP104.GNS
	54528, // MAP105.GNS
	54621, // MAP106.GNS
	54716, // MAP107.GNS
	54812, // MAP108.GNS
	54909, // MAP109.GNS
	55004, // MAP110.GNS
	55097, // MAP111.GNS
	55192, // MAP112.GNS
	55286, // MAP113.GNS
	55383, // MAP114.GNS
	56051, // MAP115.GNS
	56123, // MAP116.GNS
	56201, // MAP117.GNS
	56279, // MAP118.GNS
	56356, // MAP119.GNS
	0,     // MAP120.GNS
	0,     // MAP121.GNS
	0,     // MAP122.GNS
	0,     // MAP123.GNS
	0,     // MAP124.GNS
	56435, // MAP125.GNS
}
