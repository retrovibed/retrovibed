package tracking

import (
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/internal/md5x"
)

func HashUID(md *metainfo.Hash) string {
	return md5x.FormatUUID(md5x.Digest(md.Bytes()))
}

func PublicTrackers() []string {
	return []string{
		"udp://tracker.opentrackr.org:1337/announce",
		"udp://open.demonii.com:1337/announce",
		"udp://open.stealth.si:80/announce",
		"udp://explodie.org:6969/announce",
		"udp://tracker.torrent.eu.org:451/announce",
		"udp://exodus.desync.com:6969/announce",
		"udp://retracker01-msk-virt.corbina.net:80/announce",
		"udp://leet-tracker.moe:1337/announce",
		"udp://isk.richardsw.club:6969/announce",
		"udp://bt.ktrackers.com:6666/announce",
		"http://www.genesis-sp.org:2710/announce",
		"http://tracker.xiaoduola.xyz:6969/announce",
		"http://tracker.vanitycore.co:6969/announce",
		"http://tracker.sbsub.com:2710/announce",
		"http://tracker.moxing.party:6969/announce",
		"http://tracker.dmcomic.org:2710/announce",
		"http://tracker.bt-hash.com:80/announce",
		"http://t.jaekr.sh:6969/announce",
		"http://shubt.net:2710/announce",
		"http://share.hkg-fansub.info:80/announce.php",
	}
}
