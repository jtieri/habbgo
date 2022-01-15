package messages

import "github.com/jtieri/habbgo/protocol/packets"

func INTERSTITIALDATA() *packets.OutgoingPacket {
	p := packets.NewOutgoing(258) // Base64 Header DB

	/*
		This is in the client Lingo src:
		on handle_interstitialdata me, tMsg
		  if tMsg.content.length > 1 then
		    tDelim = the itemDelimiter
		    the itemDelimiter = "\t"
		    tSourceURL = tMsg.content.getProp(#item, 1)
		    tTargetURL = tMsg.content.getProp(#item, 2)
		    the itemDelimiter = tDelim
		    me.getComponent().getInterstitial().Init(tSourceURL, tTargetURL)
		  else
		    me.getComponent().getInterstitial().Init(0)
		  end if
		end

		Interstitial may be used to do ads on loading screens while entering rooms?
	*/

	p.WriteInt(0)
	return p
}
