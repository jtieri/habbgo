package composers

import (
	"github.com/jtieri/HabbGo/habbgo/protocol/packets"
	"github.com/jtieri/HabbGo/habbgo/utils"
	"strconv"
)


const ( // Used in ComposeSessionParams
	registerCoppa              = 0 // toggle conf_coppa or conf_strong_coppa_req by setting value > 0 or > 1
	voucherEnabled             = 1 // Enables in-game vouchers when value is set > 0
	registerRequireParentEmail = 2 // Requires parent email when registering if value is set > 0
	registerSendParentEmail    = 3 // conf_parent_email_request_reregistration
	allowDirectMail            = 4 // conf_allow_direct_mail
	dateFormat                 = 5 // Sets the date formatter used across the client
	partnerIntegrationEnabled  = 6 // conf_partner_integration. Value is either 1 or 0 (enabled or disabled)
	allowProfileEditing        = 7 // Enables the in-game profile editor
	trackingHeader             = 8 // tracking_header - used in stats.tracking.javascript(?)
	tutorialEnabled            = 9 // Enables the in-game tutorial when value is set to 1 and disables it when 0
)

func ComposeHello() *packets.OutgoingPacket {
	return packets.NewOutgoing(0) // Base64 Header @@
}

func ComposeCryptoParams() *packets.OutgoingPacket {
	packet := packets.NewOutgoing(277) // Base64 Header DU
	packet.WriteInt(0)                       // Toggles server->client encryption; 0=off | non-zero=on
	return packet
}

func ComposeEndCrypto() *packets.OutgoingPacket {
	packet := packets.NewOutgoing(278) // Base 64 Header DV
	return packet
}

func ComposeSessionParams() *packets.OutgoingPacket {
	packet := packets.NewOutgoing(257) // Base64 Header DA

	params := make(map[int]string, 10)
	params[registerCoppa] = strconv.Itoa(0)
	params[voucherEnabled] = strconv.Itoa(0) // TODO create config to enable if vouchers are enabled
	params[registerRequireParentEmail] = strconv.Itoa(0)
	params[registerSendParentEmail] = strconv.Itoa(0)
	params[allowDirectMail] = strconv.Itoa(0)
	params[dateFormat] = "dd-MM-yyyy"
	params[partnerIntegrationEnabled] = strconv.Itoa(0)
	params[allowProfileEditing] = strconv.Itoa(1) // TODO create config to enable if profile editing is enabled
	params[trackingHeader] = ""
	params[tutorialEnabled] = strconv.Itoa(0) // TODO check if player has finished tutorial then set appropriately

	packet.WriteInt(len(params))

	for i, v := range params {
		packet.WriteInt(i)

		if utils.IsNumber(v) {
			num, _ := strconv.Atoi(v)
			packet.WriteInt(num)
		} else {
			packet.WriteString(v)
		}
	}
	return packet
}

func ComposeLoginOk() *packets.OutgoingPacket {
	packet := packets.NewOutgoing(3) // Base 64 Header @C
	return packet
}
