package parser

import (
	"testing"
)
	func testHeaderParser(t *testing.T, hp HeaderParser){
			if sh, err := hp.Parse(); err!=nil{
            	t.Log(err);
            	t.Fail();
            }else{
            	t.Log("encoded = " + sh.String());
            }
	}

	func TestCallIdParser(t *testing.T) {
        var call = []string{
	    	"Call-ID :  f0b40bcc-3485-49e7-ad1a-f1dfad2e39c9@10.5.0.53\n",
            "Call-ID	: 	f81d4fae-7dec-11d0-a765-00a0c91e6bf6@foo.bar.com\n",
            "i:f81d4fae-7dec-11d0-a765-00a0c91e6bf6@foo.bar.com\n",
            "Call-ID:	1@10.0.0.1\n",
            "Call-ID	:kl24ahsd546folnyt2vbak9sad98u23naodiunzds09a3bqw0sdfbsk34poouymnae0043nsed09mfkvc74bd0cuwnms05dknw87hjpobd76f\n",
            "Call-Id: 281794\n",
        }
        
        for i := 0; i < len(call); i++  {
            shp := NewCallIDParser(call[i]);
            testHeaderParser(t, shp)
        }
    }

     func TestCSeqParser(t *testing.T) {
		var cseq =[]string{
			"CSeq: 17 INVITE\n",
			"CSeq: 17 ACK\n",
			"CSeq : 18   BYE\n",
            "CSeq:1 CANCEL\n",
            "CSeq: 3 BYE\n",
        };
			
		for i := 0; i < len(cseq); i++  {
		    shp := NewCSeqParser(cseq[i]);
		    testHeaderParser(t, shp);
		}	
	}
	
	func TestAddressParser(t *testing.T) {
	     var addresses =[]string { 
		"<sip:user@example.com?Route=%3csip:sip.example.com%3e>",
		"\"M. Ranganathan\"   <sip:mranga@nist.gov>",
		"<sip:+1-650-555-2222@ss1.wcom.com;user=phone>",
		"M. Ranganathan <sip:mranga@nist.gov>" };

	    for i := 0; i < len(addresses); i++ {
        	addressParser := NewAddressParser(addresses[i]);
			if addr, err := addressParser.Address(); err!=nil{
            	t.Log(err);
            	t.Fail();
            }else{
            	t.Log("encoded = " + addr.String());
            }
	    }
	}
