package main

import (
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/searis/guma/stack"
	"github.com/searis/guma/stack/transport"
	"github.com/searis/guma/stack/transport/tcp"
	"github.com/searis/guma/stack/uatype"
)

var typeFmt = spew.ConfigState{
	Indent:   "\t",
	SortKeys: true,
}

func main() {
	conn := tcp.Conn{}

	err := conn.Connect("localhost:4840")
	if err != nil {
		log.Fatal(err)
	}

	channel := transport.SecureChannel{}
	ec := make(chan error)
	err = channel.Open(&conn, ec)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case err := <-ec:
				log.Println(err)
			}
		}
	}()

	client := stack.Client{
		Channel: channel,
	}

	resp, err := client.CreateSession(uatype.CreateSessionRequest{
		RequestHeader: uatype.RequestHeader{
			Timestamp: time.Now(),
		},
		ClientDescription: uatype.ApplicationDescription{
			ApplicationUri: "auri",
			ProductUri:     "puri",
			ApplicationName: uatype.LocalizedText{
				TextSpecified: true,
				Text:          "Best App",
			},
			ApplicationType:     uatype.ApplicationTypeClient,
			GatewayServerUri:    "test",
			DiscoveryProfileUri: "test",
			NoOfDiscoveryUrls:   0,
			//DiscoveryUrls       []string `opcua:"lengthField=NoOfDiscoveryUrls"`
		},
		ServerUri:               "serverURI",
		EndpointUrl:             "endPointURI",
		SessionName:             "SssioName",
		ClientNonce:             uatype.ByteString(make([]byte, 32)),
		ClientCertificate:       uatype.ByteString(clientCertificate),
		RequestedSessionTimeout: 1200000,
		MaxResponseMessageSize:  16777216,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("---------- SENDING ACT SESSION ---------------")

	actresp, err := client.ActivateSession(uatype.ActivateSessionRequest{
		RequestHeader: uatype.RequestHeader{
			Timestamp:           time.Now(),
			AuthenticationToken: resp.AuthenticationToken,
		},
	})

	fmt.Printf("<< Received data:\n%s", typeFmt.Sdump(actresp))

	fmt.Println("---------- SENDING BROWSE ---------------")

	bres, err := client.Browse(uatype.BrowseRequest{
		RequestHeader: uatype.RequestHeader{
			Timestamp: time.Now(),
		},
		NodesToBrowse: []uatype.BrowseDescription{
			uatype.BrowseDescription{
				NodeId: uatype.NewFourByteNodeID(0, 84),
			},
		},
	})

	fmt.Printf("<< Received data:\n%s", typeFmt.Sdump(bres))
}

var clientCertificate = []byte{
	0x30, 0x82, 0x05, 0x05, 0x30, 0x82, 0x03, 0xed, 0xa0, 0x03, 0x02, 0x01,
	0x02, 0x02, 0x04, 0x54, 0x6c, 0x8f, 0xee, 0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86,
	0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x05, 0x05, 0x00, 0x30, 0x6c, 0x31, 0x0b,
	0x30, 0x09, 0x06, 0x03, 0x55, 0x04, 0x06, 0x13, 0x02, 0x4e, 0x4f, 0x31, 0x15,
	0x30, 0x13, 0x06, 0x03, 0x55, 0x04, 0x08, 0x13, 0x0c, 0x73, 0x6f, 0x72, 0x74,
	0x72, 0x6f, 0x6e, 0x64, 0x65, 0x6c, 0x61, 0x67, 0x31, 0x12, 0x30, 0x10, 0x06,
	0x03, 0x55, 0x04, 0x07, 0x13, 0x09, 0x74, 0x72, 0x6f, 0x6e, 0x64, 0x68, 0x65,
	0x69, 0x6d, 0x31, 0x0f, 0x30, 0x0d, 0x06, 0x03, 0x55, 0x04, 0x0a, 0x13, 0x06,
	0x73, 0x65, 0x61, 0x72, 0x69, 0x73, 0x31, 0x21, 0x30, 0x1f, 0x06, 0x03, 0x55,
	0x04, 0x03, 0x14, 0x18, 0x55, 0x61, 0x45, 0x78, 0x70, 0x65, 0x72, 0x74, 0x40,
	0x42, 0x45, 0x52, 0x4e, 0x54, 0x2d, 0x4a, 0x4f, 0x48, 0x41, 0x4e, 0x42, 0x34,
	0x39, 0x44, 0x30, 0x1e, 0x17, 0x0d, 0x31, 0x34, 0x31, 0x31, 0x31, 0x39, 0x31,
	0x32, 0x34, 0x31, 0x31, 0x38, 0x5a, 0x17, 0x0d, 0x31, 0x35, 0x31, 0x31, 0x31,
	0x39, 0x31, 0x32, 0x34, 0x31, 0x31, 0x38, 0x5a, 0x30, 0x6c, 0x31, 0x0b, 0x30,
	0x09, 0x06, 0x03, 0x55, 0x04, 0x06, 0x13, 0x02, 0x4e, 0x4f, 0x31, 0x15, 0x30,
	0x13, 0x06, 0x03, 0x55, 0x04, 0x08, 0x13, 0x0c, 0x73, 0x6f, 0x72, 0x74, 0x72,
	0x6f, 0x6e, 0x64, 0x65, 0x6c, 0x61, 0x67, 0x31, 0x12, 0x30, 0x10, 0x06, 0x03,
	0x55, 0x04, 0x07, 0x13, 0x09, 0x74, 0x72, 0x6f, 0x6e, 0x64, 0x68, 0x65, 0x69,
	0x6d, 0x31, 0x0f, 0x30, 0x0d, 0x06, 0x03, 0x55, 0x04, 0x0a, 0x13, 0x06, 0x73,
	0x65, 0x61, 0x72, 0x69, 0x73, 0x31, 0x21, 0x30, 0x1f, 0x06, 0x03, 0x55, 0x04,
	0x03, 0x14, 0x18, 0x55, 0x61, 0x45, 0x78, 0x70, 0x65, 0x72, 0x74, 0x40, 0x42,
	0x45, 0x52, 0x4e, 0x54, 0x2d, 0x4a, 0x4f, 0x48, 0x41, 0x4e, 0x42, 0x34, 0x39,
	0x44, 0x30, 0x82, 0x01, 0x22, 0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86,
	0xf7, 0x0d, 0x01, 0x01, 0x01, 0x05, 0x00, 0x03, 0x82, 0x01, 0x0f, 0x00, 0x30,
	0x82, 0x01, 0x0a, 0x02, 0x82, 0x01, 0x01, 0x00, 0x97, 0xc9, 0xb2, 0xde, 0x50,
	0xbb, 0x3e, 0x9d, 0x24, 0x81, 0x58, 0xc3, 0x75, 0xad, 0xa9, 0xaa, 0x7f, 0x7e,
	0xcf, 0x9b, 0x59, 0x43, 0xdc, 0xfc, 0xc0, 0x09, 0x91, 0x93, 0xda, 0xad, 0x09,
	0x82, 0xb7, 0x2f, 0x1c, 0xf1, 0x8c, 0x5d, 0xd1, 0x26, 0xf5, 0x0f, 0x93, 0xec,
	0xc8, 0x8f, 0xd9, 0x37, 0x79, 0xa5, 0x27, 0x10, 0x6a, 0x85, 0xcf, 0x75, 0xd5,
	0x9f, 0xd0, 0x52, 0xad, 0x69, 0xba, 0x47, 0xc8, 0x8b, 0xe4, 0xbf, 0x09, 0x9e,
	0xd5, 0x19, 0x5b, 0x41, 0xb1, 0x8f, 0x32, 0x09, 0x32, 0x9f, 0x16, 0xb7, 0x98,
	0x1b, 0x7d, 0xe5, 0x2e, 0x45, 0x70, 0xe7, 0xd5, 0x55, 0xf2, 0x22, 0x7c, 0xb9,
	0xb6, 0xd6, 0x8d, 0xc0, 0x7c, 0x85, 0xb3, 0xa4, 0xe6, 0x47, 0x11, 0x6d, 0x9e,
	0xf1, 0xb3, 0x69, 0x91, 0x42, 0xf2, 0xae, 0x8c, 0xdf, 0x9b, 0xdf, 0x67, 0xbb,
	0xcd, 0x11, 0x62, 0x2e, 0x81, 0x62, 0xfa, 0xe0, 0xc4, 0x4b, 0xbd, 0xe4, 0x59,
	0x79, 0x7a, 0xf7, 0x39, 0x11, 0xc2, 0x5a, 0xb7, 0x78, 0xb7, 0xe1, 0x28, 0x0e,
	0xdd, 0x09, 0xcc, 0x71, 0x2f, 0xcf, 0xc3, 0xb2, 0x5a, 0xf4, 0x0a, 0x80, 0xe7,
	0x70, 0xf1, 0x92, 0x33, 0xde, 0x09, 0xa7, 0x50, 0x03, 0x1c, 0x52, 0xbf, 0xde,
	0xfc, 0xaa, 0x9e, 0x86, 0x2f, 0x69, 0x74, 0x3f, 0xec, 0x17, 0x88, 0x70, 0xa2,
	0xfd, 0xf2, 0x39, 0xca, 0xfb, 0x04, 0x1f, 0x2d, 0x63, 0xc2, 0x7c, 0x3a, 0x62,
	0xa9, 0xae, 0xb8, 0x32, 0x0c, 0x50, 0x60, 0x2a, 0x4f, 0x33, 0x70, 0x3d, 0x16,
	0x27, 0x03, 0xa4, 0xca, 0x23, 0x37, 0x9b, 0x42, 0x35, 0x81, 0x28, 0x90, 0x62,
	0xb6, 0x49, 0x9e, 0x7b, 0x5f, 0xbb, 0xec, 0xbf, 0xee, 0xef, 0xc5, 0x86, 0x9c,
	0x45, 0x3b, 0x0d, 0x4b, 0xdb, 0x03, 0xe4, 0x31, 0x3f, 0x6c, 0x2c, 0x1b, 0x60,
	0xeb, 0xf8, 0xd8, 0x09, 0x02, 0x03, 0x01, 0x00, 0x01, 0xa3, 0x82, 0x01, 0xad,
	0x30, 0x82, 0x01, 0xa9, 0x30, 0x0c, 0x06, 0x03, 0x55, 0x1d, 0x13, 0x01, 0x01,
	0xff, 0x04, 0x02, 0x30, 0x00, 0x30, 0x50, 0x06, 0x09, 0x60, 0x86, 0x48, 0x01,
	0x86, 0xf8, 0x42, 0x01, 0x0d, 0x04, 0x43, 0x16, 0x41, 0x22, 0x47, 0x65, 0x6e,
	0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x55,
	0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x20, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x20, 0x55, 0x41, 0x20, 0x42, 0x61, 0x73, 0x65, 0x20,
	0x4c, 0x69, 0x62, 0x72, 0x61, 0x72, 0x79, 0x20, 0x75, 0x73, 0x69, 0x6e, 0x67,
	0x20, 0x4f, 0x70, 0x65, 0x6e, 0x53, 0x53, 0x4c, 0x22, 0x30, 0x1d, 0x06, 0x03,
	0x55, 0x1d, 0x0e, 0x04, 0x16, 0x04, 0x14, 0x3b, 0xbb, 0x14, 0x58, 0xf8, 0x5b,
	0xe4, 0x58, 0x79, 0x4f, 0xfd, 0x5e, 0x82, 0xfd, 0xc2, 0x11, 0x33, 0x5c, 0x6a,
	0x89, 0x30, 0x81, 0x99, 0x06, 0x03, 0x55, 0x1d, 0x23, 0x04, 0x81, 0x91, 0x30,
	0x81, 0x8e, 0x80, 0x14, 0x3b, 0xbb, 0x14, 0x58, 0xf8, 0x5b, 0xe4, 0x58, 0x79,
	0x4f, 0xfd, 0x5e, 0x82, 0xfd, 0xc2, 0x11, 0x33, 0x5c, 0x6a, 0x89, 0xa1, 0x70,
	0xa4, 0x6e, 0x30, 0x6c, 0x31, 0x0b, 0x30, 0x09, 0x06, 0x03, 0x55, 0x04, 0x06,
	0x13, 0x02, 0x4e, 0x4f, 0x31, 0x15, 0x30, 0x13, 0x06, 0x03, 0x55, 0x04, 0x08,
	0x13, 0x0c, 0x73, 0x6f, 0x72, 0x74, 0x72, 0x6f, 0x6e, 0x64, 0x65, 0x6c, 0x61,
	0x67, 0x31, 0x12, 0x30, 0x10, 0x06, 0x03, 0x55, 0x04, 0x07, 0x13, 0x09, 0x74,
	0x72, 0x6f, 0x6e, 0x64, 0x68, 0x65, 0x69, 0x6d, 0x31, 0x0f, 0x30, 0x0d, 0x06,
	0x03, 0x55, 0x04, 0x0a, 0x13, 0x06, 0x73, 0x65, 0x61, 0x72, 0x69, 0x73, 0x31,
	0x21, 0x30, 0x1f, 0x06, 0x03, 0x55, 0x04, 0x03, 0x14, 0x18, 0x55, 0x61, 0x45,
	0x78, 0x70, 0x65, 0x72, 0x74, 0x40, 0x42, 0x45, 0x52, 0x4e, 0x54, 0x2d, 0x4a,
	0x4f, 0x48, 0x41, 0x4e, 0x42, 0x34, 0x39, 0x44, 0x82, 0x04, 0x54, 0x6c, 0x8f,
	0xee, 0x30, 0x0e, 0x06, 0x03, 0x55, 0x1d, 0x0f, 0x01, 0x01, 0xff, 0x04, 0x04,
	0x03, 0x02, 0x02, 0xf4, 0x30, 0x20, 0x06, 0x03, 0x55, 0x1d, 0x25, 0x01, 0x01,
	0xff, 0x04, 0x16, 0x30, 0x14, 0x06, 0x08, 0x2b, 0x06, 0x01, 0x05, 0x05, 0x07,
	0x03, 0x01, 0x06, 0x08, 0x2b, 0x06, 0x01, 0x05, 0x05, 0x07, 0x03, 0x02, 0x30,
	0x5a, 0x06, 0x03, 0x55, 0x1d, 0x11, 0x04, 0x53, 0x30, 0x51, 0x86, 0x3e, 0x75,
	0x72, 0x6e, 0x3a, 0x42, 0x45, 0x52, 0x4e, 0x54, 0x2d, 0x4a, 0x4f, 0x48, 0x41,
	0x4e, 0x42, 0x34, 0x39, 0x44, 0x3a, 0x55, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64,
	0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x55, 0x61,
	0x45, 0x78, 0x70, 0x65, 0x72, 0x74, 0x40, 0x42, 0x45, 0x52, 0x4e, 0x54, 0x2d,
	0x4a, 0x4f, 0x48, 0x41, 0x4e, 0x42, 0x34, 0x39, 0x44, 0x82, 0x0f, 0x42, 0x45,
	0x52, 0x4e, 0x54, 0x2d, 0x4a, 0x4f, 0x48, 0x41, 0x4e, 0x42, 0x34, 0x39, 0x44,
	0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x05,
	0x05, 0x00, 0x03, 0x82, 0x01, 0x01, 0x00, 0x12, 0x8e, 0xd7, 0xfc, 0x01, 0xdf,
	0x7a, 0x38, 0xc8, 0x4c, 0x38, 0xe0, 0xab, 0x89, 0xcb, 0x64, 0x86, 0x11, 0xd8,
	0xbe, 0xaa, 0xd3, 0x3a, 0xa3, 0x26, 0x8f, 0x82, 0xb3, 0xbe, 0x89, 0xc3, 0x1b,
	0x77, 0x29, 0x5f, 0x22, 0xed, 0x8e, 0xfc, 0xfb, 0xa8, 0x0d, 0x06, 0xb2, 0x40,
	0x2d, 0x0f, 0x67, 0x56, 0xe4, 0x69, 0x01, 0x51, 0xef, 0xb5, 0x22, 0xd4, 0x05,
	0x79, 0x7c, 0x86, 0xe6, 0x7c, 0x08, 0x9a, 0xc4, 0xf5, 0x17, 0x20, 0x1b, 0xca,
	0xec, 0xb3, 0x06, 0x76, 0xc4, 0x36, 0xe0, 0x5c, 0x97, 0xd5, 0xf7, 0x1d, 0x3f,
	0xbf, 0xa9, 0x6a, 0x7f, 0xa6, 0x67, 0x4e, 0x6d, 0xce, 0x0c, 0x3e, 0xae, 0xbf,
	0x7d, 0x1d, 0xb8, 0xef, 0x86, 0x5e, 0xc3, 0xfb, 0xc6, 0x6b, 0x2b, 0x0b, 0x0c,
	0xe3, 0x4b, 0xd9, 0x28, 0xe5, 0x65, 0x3f, 0xac, 0x3d, 0x88, 0xb0, 0x5f, 0x1d,
	0x3c, 0x1a, 0x72, 0x23, 0x95, 0xd5, 0xd1, 0x76, 0x1b, 0xfb, 0x59, 0x32, 0x53,
	0x3b, 0x77, 0x0a, 0xcf, 0xb5, 0x7e, 0xab, 0xf2, 0xa0, 0x64, 0x00, 0x66, 0x64,
	0xf3, 0x4f, 0x2b, 0x7f, 0x17, 0xfe, 0x94, 0x4f, 0xbe, 0x6a, 0x58, 0x25, 0x17,
	0x9a, 0x5a, 0xdc, 0x20, 0x1c, 0x30, 0x3e, 0x81, 0xcf, 0x92, 0x4d, 0x5d, 0xf7,
	0x10, 0x96, 0x39, 0xc6, 0xb7, 0x43, 0x83, 0xb4, 0x29, 0x16, 0x03, 0x66, 0x6e,
	0xb3, 0x1f, 0x4e, 0x21, 0xa6, 0x60, 0x37, 0xe7, 0xe6, 0xcb, 0xbb, 0x7f, 0x96,
	0xb6, 0x10, 0xc0, 0x69, 0x45, 0x8f, 0xe0, 0xba, 0x13, 0x5d, 0x26, 0x5a, 0x90,
	0x20, 0x53, 0xfb, 0x69, 0x91, 0x89, 0xa3, 0xfc, 0x01, 0xb0, 0xe8, 0x83, 0x5f,
	0x30, 0xfe, 0x6d, 0xf4, 0xb6, 0x69, 0x28, 0x67, 0xdf, 0xf0, 0xbe, 0x75, 0xdf,
	0x05, 0x65, 0xf1, 0x24, 0x2e, 0x4d, 0x6d, 0x7d, 0x9b, 0x82, 0x5e, 0xa7, 0x7d,
	0x79, 0xf5, 0x5d}
