package parser

import (
	//"bytes"
	"errors"
	"gosip/address"
	"gosip/core"
	"gosip/header"
	"gosip/message"
	"strings"
)

/**
 * Parse SIP message and parts of SIP messages such as URI's etc
 * from memory and return a structure.
 * Intended use:  UDP message processing.
 * This class is used when you have an entire SIP message or SIPHeader
 * or SIP URL in memory and you want to generate a parsed structure from
 * it. For SIP messages, the payload can be binary or String.
 * If you have a binary payload,
 * use parseSIPMessage(byte[]) else use parseSIPMessage(String)
 * The payload is accessible from the parsed message using the GetContent and
 * GetContentBytes methods provided by the SIPMessage class. If SDP parsing
 * is enabled using the parseContent method, then the SDP body is also parsed
 * and can be accessed from the message using the GetSDPAnnounce method.
 * Currently only eager parsing of the message is supported (i.e. the
 * entire message is parsed in one feld swoop).
 *
 *
 *@version  JAIN-SIP-1.1
 *
 *@author M. Ranganathan <mranga@nist.gov>  <br/>
 *
 *<a href="{@docRoot}/uncopyright.html">This code is in the public domain.</a>
 *
 */

type StringMsgParser struct {
	readBody bool

	rawMessage string
	// Unprocessed message  (for error reporting)
	rawMessage1 string
	// Unprocessed message  (for error reporting)
	currentMessage string
	// the message being parsed. (for error reporting)
	//private ParseExceptionListener parseExceptionListener;

	messageHeaders map[int]string //header.ISIPHeader // Message headers

	bufferPointer int

	bodyIsString bool

	currentMessageBytes []byte

	contentLength int

	debugFlag bool

	currentLine int

	currentHeader string
}

/**
 *@since v0.9
 */
func NewStringMsgParser() *StringMsgParser {
	this := &StringMsgParser{}
	//super();
	//this.messageHeaders = new Vector(10,10);
	this.bufferPointer = 0
	this.currentLine = 0
	this.readBody = true
	return this
}

/**
 *Constructor (given a parse exception handler).
 *@since 1.0
 *@param exhandler is the parse exception listener for the message parser.
 */
// func NewStringMsgParser( ParseExceptionListener exhandler) {
//     this();
//     parseExceptionListener = exhandler;
// }

/** Get the message body.
 */
func (this *StringMsgParser) GetMessageBody() string {
	if this.contentLength == 0 {
		return ""
	} else {
		endIndex := this.bufferPointer + this.contentLength
		var body string
		// guard against bad specifications.
		if endIndex > len(this.currentMessage) {
			endIndex = len(this.currentMessage)
			body = this.currentMessage[this.bufferPointer:endIndex]
			this.bufferPointer = endIndex
		} else {
			body = this.currentMessage[this.bufferPointer:endIndex]
			this.bufferPointer = endIndex + 1
		}
		this.contentLength = 0
		return body
	}

}

/** Get the message body as a byte array.
 */
func (this *StringMsgParser) GetBodyAsBytes() []byte {
	if this.contentLength == 0 {
		return nil
	} else {
		endIndex := this.bufferPointer + this.contentLength
		// guard against bad specifications.
		if endIndex > len(this.currentMessageBytes) {
			endIndex = len(this.currentMessageBytes)
		}
		body := make([]byte, endIndex-this.bufferPointer)

		copy(body, this.currentMessageBytes[this.bufferPointer:endIndex])
		//
		//            for (int i = bufferPointer, k = 0; i < endIndex; i++,k++) {
		//                body[k] = currentMessageBytes[i];
		//            }

		this.bufferPointer = endIndex
		this.contentLength = 0
		return body
	}

}

/** Return the contents till the end of the buffer (this is useful when
 * you encounter an error.
 */
func (this *StringMsgParser) ReadToEnd() string {
	body := this.currentMessage[this.bufferPointer:]
	this.bufferPointer += len(body)
	return body
}

/** Return tbe bytes to the end of the message.
 * This is invoked when the parser is invoked with an array of bytes
 * rather than with a string.
 */
func (this *StringMsgParser) ReadBytesToEnd() []byte {
	body := make([]byte, len(this.currentMessageBytes)-this.bufferPointer)
	endIndex := len(this.currentMessageBytes)
	copy(body, this.currentMessageBytes[this.bufferPointer:endIndex])
	this.bufferPointer = endIndex
	this.contentLength = 0
	return body
}

/**
 * add a handler for header parsing errors.
 * @param  pexhadler is a class
 *  	that implements the ParseExceptionListener interface.
 */

// public void SetParseExceptionListener
// ( ParseExceptionListener pexhandler ) {
//     parseExceptionListener = pexhandler;
// }

/** Return true if the body is encoded as a string.
 * If the parseSIPMessage(String) method is invoked then the body
 * is assumed to be a string.
 */
func (this *StringMsgParser) IsBodyString() bool {
	return this.bodyIsString
}

/** Parse a buffer containing a single SIP Message where the body
 * is an array of un-interpreted bytes. This is intended for parsing
 * the message from a memory buffer when the buffer.
 * Incorporates a bug fix for a bug that was noted by Will Sullin of
 * Callcast
 * @param msgBuffer a byte buffer containing the messages to be parsed.
 *   This can consist of multiple SIP Messages concatenated toGether.
 * @return a SIPMessage[] structure (request or response)
 * 			containing the parsed SIP message.
 * @exception SIPIllegalMessageException is thrown when an
 * 			illegal message has been encountered (and
 *			the rest of the buffer is discarded).
 * @see ParseExceptionListener
 */
// func (this *StringMsgParser) ParseSIPMessage(msgBuffer []byte) (*message.SIPMessage, error) {
// 	//throws ParseException {
// 	this.bufferPointer = 0
// 	this.bodyIsString = false
// 	//retval := new Vector();
// 	this.currentMessageBytes = msgBuffer
// 	var s int
// 	// Squeeze out leading CRLF
// 	// Squeeze out the leading nulls (otherwise the parser will crash)
// 	// Bug noted by Will Sullin of Callcast
// 	for s = this.bufferPointer; s < len(msgBuffer); s++ {
// 		if msgBuffer[s] != '\r' && msgBuffer[s] != '\n' {
// 			// msgBuffer[s] != '\0' {
// 			break
// 		}
// 	}

// 	if s == len(msgBuffer) {
// 		return nil, nil
// 	}

// 	// Find the end of the SIP message.
// 	var f int
// 	for f = s; f < len(msgBuffer)-4; f++ {
// 		if msgBuffer[f] == '\r' &&
// 			msgBuffer[f+1] == '\n' &&
// 			msgBuffer[f+2] == '\r' &&
// 			msgBuffer[f+3] == '\n' {
// 			break
// 		}
// 	}
// 	if f < len(msgBuffer) {
// 		f += 4
// 	} else {
// 		// Could not find CRLFCRLF end of message so look for LFLF
// 		for f = s; f < len(msgBuffer)-2; f++ {
// 			if msgBuffer[f] == '\n' &&
// 				msgBuffer[f] == '\n' {
// 				break
// 			}
// 		}
// 		if f < len(msgBuffer) {
// 			f += 2
// 		} else {
// 			return nil, errors.New("ParseException: Message not terminated")
// 		}
// 	}

// 	// Encode the body as a UTF-8 string.
// 	var messageString string
// 	//try {
// 	messageString = string(msgBuffer[s : f-s]) //, "UTF-8");
// 	//
// 	this.bufferPointer = f
// 	message := []byte(messageString)
// 	length := len(message)
// 	// Get rid of CR to make it uniform for the parser.
// 	for k := 0; k < length; k++ {
// 		if message[k] == '\r' {
// 			copy(message[k:length-1], message[k+1:length])
// 			length--
// 		}
// 	}

// 	// if (Parser.debug) {
// 	//     for (int k = 0 ; k < length; k++) {
// 	//         rawMessage1 = rawMessage1 + "[" + message.charAt(k) +"]";
// 	//     }
// 	// }

// 	// The following can be written more efficiently in a single pass
// 	// but it is somewhat tricky.
// 	tokenizer := core.NewStringTokenizer(string(message[:length])) //,"\n",true);
// 	var cooked_message bytes.Buffer                                //= new StringBuffer();
// 	//try {
// 	for tokenizer.HasMoreChars() { //hasMoreElements() ) {
// 		nexttok := tokenizer.GetNextTokenByDelim('\n') //nextToken();
// 		// Ignore blank lines with leading spaces or tabs.
// 		if strings.TrimSpace(nexttok) == "" {
// 			cooked_message.WriteString("\n")
// 		} else {
// 			cooked_message.WriteString(nexttok)
// 		}
// 	}
// 	// } catch (NoSuchElementException ex) {
// 	// }

// 	message1 := cooked_message.String()
// 	length = message1.indexOf("\n\n") + 2

// 	// Handle continuations - look for a space or a tab at the start
// 	// of the line and append it to the previous line.

// 	for k := 0; k < length-1; {
// 		if cooked_message[k] == '\n' {
// 			if cooked_message[k+1] == '\t' || cooked_message[k+1] == ' ' {
// 				cooked_message.deleteCharAt(k)
// 				cooked_message.deleteCharAt(k)
// 				length--
// 				length--
// 				if k == length {
// 					break
// 				} else {
// 					continue
// 				}
// 			}

// 			if cooked_message[k+1] == '\n' {
// 				cooked_message.insert(k, '\n')
// 				length++
// 				k++
// 			}
// 		}
// 		k++
// 	}
// 	cooked_message.WriteString("\n\n")

// 	// Separate the string out into substrings for
// 	// error reporting.
// 	this.currentMessage = cooked_message.String()
// 	sipmsg := this.ParseMessage(this.currentMessage)
// 	if this.readBody && sipmsg.GetContentLength() != null && sipmsg.GetContentLength().GetContentLength() != 0 {
// 		this.contentLength = sipmsg.GetContentLength().GetContentLength()
// 		body := this.GetBodyAsBytes()
// 		sipmsg.SetMessageContent(body)
// 	}
// 	// System.out.println("Parsed = " + sipmsg);
// 	return sipmsg, nil
// }

//     /**
//      * Parse a buffer containing one or more SIP Messages  and return an array of
//      * SIPMessage parsed structures. Note that the current limitation is that
//      * this does not handle content encoding properly. The message content is
//      * just assumed to be encoded using the same encoding as the sip message
//      * itself (i.e. binary encodings such as gzip are not supported).
//      * @param sipMessages a String containing the messages to be parsed.
//      *   This can consist of multiple SIP Messages concatenated toGether.
//      * @return a SIPMessage structure (request or response)
//      * 			containing the parsed SIP message.
//      * @exception SIPIllegalMessageException is thrown when an
//      * 			illegal message has been encountered (and
//      *			the rest of the buffer is discarded).
//      * @see ParseExceptionListener
//      */

//     public SIPMessage  parseSIPMessage(String sipMessages )
//     throws ParseException {
//         // Handle line folding and evil DOS CR-LF sequences
//         rawMessage = sipMessages;
//         Vector retval = new Vector();
//         String pmessage = sipMessages.trim();
//         bodyIsString = true;

//             this.contentLength = 0;
//             if (pmessage.trim().equals("")) return null;

//             pmessage += "\n\n";
//             StringBuffer message = new StringBuffer(pmessage);
//             // squeeze out the leading crlf sequences.
//             while(message.charAt(0) == '\r' || message.charAt(0) == '\n') {
//                 bufferPointer ++;
//                 message.deleteCharAt(0);
//             }

//             // squeeze out the crlf sequences and make them uniformly CR
//             String message1 = message.toString();
//             int length;
//             length = message1.indexOf("\r\n\r\n");
//             if (length > 0 ) length += 4;
//             if (length == -1) {
//                 length = message1.indexOf("\n\n");
//                 if (length == -1)
//                     throw new ParseException("no trailing crlf",0);
//             } else length += 2;

//             // Get rid of CR to make it uniform.
//             for ( int k = 0; k < length ; k++ ) {
//                 if (message.charAt(k) == '\r' ) {
//                     message.deleteCharAt(k);
//                     length --;
//                 }
//             }

//             if (debugFlag ) {
//                 for (int k = 0 ; k < length; k++) {
//                     rawMessage1 = rawMessage1 + "[" + message.charAt(k) +"]";
//                 }
//             }

//             // The following can be written more efficiently in a single pass
//             // but it is somewhat tricky.
//             java.util.StringTokenizer tokenizer = new java.util.StringTokenizer
//             (message.toString(),"\n",true);
//             StringBuffer cooked_message = new StringBuffer();
//             try {
//                 while( tokenizer.hasMoreElements() ) {
//                     String nexttok = tokenizer.nextToken();
//                     // Ignore blank lines with leading spaces or tabs.
//                     if (nexttok.trim().equals("")) cooked_message.append("\n");
//                     else cooked_message.append(nexttok);
//                 }
//             } catch (NoSuchElementException ex) {
//             }

//             message1 = cooked_message.toString();
//             length = message1.indexOf("\n\n") + 2;

//             // Handle continuations - look for a space or a tab at the start
//             // of the line and append it to the previous line.

//             for ( int k = 0 ; k < length - 1 ;  ) {
//                 if (cooked_message.charAt(k) == '\n') {
//                     if ( cooked_message.charAt(k+1) == '\t' ||
//                     cooked_message.charAt(k+1) == ' ') {
//                         cooked_message.deleteCharAt(k);
//                         cooked_message.deleteCharAt(k);
//                         length --;
//                         length --;
//                         if ( k == length) break;
//                         continue;
//                     }
//                     if ( cooked_message.charAt(k+1) == '\n') {
//                         cooked_message.insert(k,'\n');
//                         length ++;
//                         k ++;
//                     }
//                 }
//                 k++;
//             }
//             cooked_message.append("\n\n");

//             // Separate the string out into substrings for
//             // error reporting.

//             currentMessage = cooked_message.toString();
//             if (Parser.debug) Debug.println(currentMessage);
//             bufferPointer = currentMessage.indexOf("\n\n") + 3 ;
//             SIPMessage sipmsg = this.parseMessage(currentMessage);
// 	    if (readBody && sipmsg.GetContentLength() != null &&
// 		sipmsg.GetContentLength().GetContentLength() != 0) {
// 		this.contentLength =
// 		     sipmsg.GetContentLength().GetContentLength();
// 		String body = this.GetMessageBody();
// 		sipmsg.SetMessageContent(body);
// 	    }
// 	    return sipmsg;

//     }

/** This is called repeatedly by parseSIPMessage to parse
 * the contents of a message buffer. This assumes the message
 * already has continuations etc. taken care of.
 * prior to its being called.
 */
func (this *StringMsgParser) ParseMessage(currentMessage string) (message.Message, error) {
	//throws  ParseException {
	// position line counter at the end of the
	// sip messages.
	// System.out.println("parsing " + currentMessage);

	sip_message_size := 0 // # of lines in the sip message
	var sipmsg message.Message

	tokenizer := core.NewStringTokenizer(currentMessage) //,"\n",true);
	this.messageHeaders = make(map[int]string)           // new Vector(); // A list of headers for error reporting
	//try {
	for tokenizer.HasMoreChars() { //hasMoreElements() ) {
		nexttok := tokenizer.GetNextTokenByDelim('\n') //nextToken();
		if nexttok == "\n" {
			nextnexttok := tokenizer.GetNextTokenByDelim('\n') //nextToken();
			if nextnexttok == "\n" {
				break
			} else {
				this.messageHeaders[sip_message_size] = nextnexttok
			}
		} else {
			this.messageHeaders[sip_message_size] = nexttok
		}
		sip_message_size++
	}
	// } catch (NoSuchElementException ex) {
	// }
	this.currentLine = 0
	this.currentHeader = this.messageHeaders[this.currentLine]
	firstLine := this.currentHeader
	// System.out.println("first Line " + firstLine);

	if !strings.Contains(firstLine, header.SIPConstants_SIP_VERSION_STRING) {
		sipmsg = message.NewSIPRequest()
		//try {
		rl, _ := NewRequestLineParser(firstLine + "\n").Parse()
		sipmsg.(*message.SIPRequest).SetRequestLine(rl)
		// } catch (ParseException ex) {
		//         if (this.parseExceptionListener != null)
		//             this.parseExceptionListener.handleException
		//             (ex,sipmsg, RequestLine.class,
		//             firstLine,currentMessage);
		//         else throw ex;

		// }
	} else {
		sipmsg = message.NewSIPResponse()
		//try {
		sl, _ := NewStatusLineParser(firstLine + "\n").Parse()
		sipmsg.(*message.SIPResponse).SetStatusLine(sl)
		// } catch (ParseException ex) {
		//         if (this.parseExceptionListener != null)   {
		//             this.parseExceptionListener.handleException
		//             (ex,sipmsg,
		//             StatusLine.class,
		//             firstLine,currentMessage);
		//         } else throw ex;

		// }
	}

	for i := 1; i < len(this.messageHeaders); i++ {
		hdrstring := this.messageHeaders[i]
		if hdrstring == "" || strings.TrimSpace(hdrstring) == "" {
			continue
		}
		///try {
		hdrParser := CreateParser(hdrstring + "\n")
		//           } catch (ParseException ex)  {
		//               this.parseExceptionListener.handleException
		//                 ( ex, sipmsg,  null , hdrstring,currentMessage);
		// continue;
		//    }
		//try {
		sipHeader, _ := hdrParser.Parse()
		if _, ok := sipmsg.(*message.SIPRequest); ok {
			sipmsg.(*message.SIPRequest).AttachHeader2(sipHeader, false)
		} else {
			sipmsg.(*message.SIPResponse).AttachHeader2(sipHeader, false)
		}
		// } catch (ParseException ex) {
		//     if (this.parseExceptionListener != null) {
		//         String hdrName = Lexer.GetHeaderName(hdrstring);
		//         Class hdrClass = NameMap.GetClassFromName(hdrName);
		//         try {
		//             if (hdrClass == null) {
		//                 hdrClass = Class.forName
		//                 (PackageNames.SIPHEADERS_PACKAGE
		//                 + ".ExtensionHeaderImpl");
		//             }
		//             this.parseExceptionListener.handleException
		//             (ex,sipmsg, hdrClass,
		//             hdrstring,currentMessage);
		//         } catch (ClassNotFoundException ex1) {
		//             InternalErrorHandler.handleException(ex1);
		//         }
		//     }

		// }
	}
	return sipmsg, nil
}

/**
 * Parse an address (nameaddr or address spec)  and return and address
 * structure.
 * @param address is a String containing the address to be parsed.
 * @return a parsed address structure.
 * @since v1.0
 * @exception  ParseException when the address is badly formatted.
 */

func (this *StringMsgParser) ParseAddress(address string) (*address.AddressImpl, error) {
	//throws ParseException {
	addressParser := NewAddressParser(address)
	return addressParser.Address()
}

/**
 * Parse a host:port and return a parsed structure.
 * @param hostport is a String containing the host:port to be parsed
 * @return a parsed address structure.
 * @since v1.0
 * @exception throws a ParseException when the address is badly formatted.
 */
func (this *StringMsgParser) ParseHostPort(hostport string) (*core.HostPort, error) {
	//throws ParseException {
	lexer := NewLexer("charLexer", hostport)
	return core.NewHostNameParserFromLexer(lexer).GetHostPort()
}

/**
 * Parse a host name and return a parsed structure.
 * @param host is a String containing the host name to be parsed
 * @return a parsed address structure.
 * @since v1.0
 * @exception throws a ParseException when the hostname is badly formatted.
 */
func (this *StringMsgParser) ParseHost(host string) (*core.Host, error) {
	//throws ParseException {
	lexer := NewLexer("charLexer", host)
	return core.NewHostNameParserFromLexer(lexer).GetHost()
}

/**
 * Parse a telephone number return a parsed structure.
 * @param telphone_number is a String containing
 * the telephone # to be parsed
 * @return a parsed address structure.
 * @since v1.0
 * @exception throws a ParseException when the address is badly formatted.
 */
func (this *StringMsgParser) ParSetelephoneNumber(telephone_number string) (*core.TelephoneNumber, error) {
	//throws ParseException {
	// Bug fix contributed by Will Scullin
	return NewURLParser(telephone_number).ParseTelephoneNumber()
}

/**
 * Parse a  SIP url from a string and return a URI structure for it.
 * @param url a String containing the URI structure to be parsed.
 * @return A parsed URI structure
 * @exception ParseException  if there was an error parsing the message.
 */
func (this *StringMsgParser) ParseSIPUrl(url string) (*address.SipUri, error) {
	//throws ParseException {
	//try {
	sipuri, err := NewURLParser(url).Parse()
	return sipuri.(*address.SipUri), err
	//} catch (ClassCastException ex) {
	//    throw new ParseException(url + " Not a SIP URL ",0);
	//}
}

/**
 * Parse a  uri from a string and return a URI structure for it.
 * @param url a String containing the URI structure to be parsed.
 * @return A parsed URI structure
 * @exception ParseException  if there was an error parsing the message.
 */
func (this *StringMsgParser) ParseUrl(url string) (*address.GenericURI, error) {
	//throws ParseException {
	guri, err := NewURLParser(url).Parse()
	return guri.(*address.GenericURI), err
}

/**
 * Parse an individual SIP message header from a string.
 * @param header String containing the SIP header.
 * @return a SIPHeader structure.
 * @exception ParseException  if there was an error parsing the message.
 */
func (this *StringMsgParser) ParseSIPHeader(h string) (*header.SIPHeader, error) {
	//throws ParseException {
	h += "\n\n"
	// Handle line folding.
	var nmessage string
	//counter := 0
	// eat leading spaces and carriage returns (necessary??)
	i := 0
	for h[i] == '\n' || h[i] == '\t' || h[i] == ' ' {
		i++
	}
	for ; i < len(h); i++ {
		if i < len(h)-1 && (h[i] == '\n' && (h[i+1] == '\t' || h[i+1] == ' ')) {
			nmessage += " "
			i++
		} else {
			nmessage += string(h[i])
		}
	}

	nmessage += "\n"

	hp := CreateParser(nmessage)
	if hp == nil {
		return nil, errors.New("ParseException: could not create parser")
	}

	shp, err := hp.Parse()
	return shp.(*header.SIPHeader), err
}

/**
 * Parse the SIP Request Line
 * @param  requestLine a String  containing the request line to be parsed.
 * @return  a RequestLine structure that has the parsed RequestLine
 * @exception ParseException  if there was an error parsing the requestLine.
 */

func (this *StringMsgParser) ParseSIPRequestLine(requestLine string) (*header.RequestLine, error) {
	//throws ParseException {
	requestLine += "\n"
	return NewRequestLineParser(requestLine).Parse()
}

/**
 * Parse the SIP Response message status line
 * @param statusLine a String containing the Status line to be parsed.
 * @return StatusLine class corresponding to message
 * @exception ParseException  if there was an error parsing
 * @see StatusLine
 */

func (this *StringMsgParser) ParseSIPStatusLine(statusLine string) (*header.StatusLine, error) {
	//throws ParseException {
	statusLine += "\n"
	return NewStatusLineParser(statusLine).Parse()
}

/**
 * Get the current header.
 */
func (this *StringMsgParser) GetCurrentHeader() string {
	return this.currentHeader
}

//     /**
//      * Get the current line number.
//      */
//     public int GetCurrentLineNumber() { return currentLine; }

/*
    public static void main(String[] args) throws ParseException {
        String messages[] = {
"SIP/2.0 180 Ringing\r\n"+
"Via: SIP/2.0/UDP 172.18.1.29:5060;branch=z9hG4bK43fc10fb4446d55fc5c8f969607991f4\r\n"+
"To: \"0440\" <sip:0440@212.209.220.131>;tag=2600\r\n"+
"From: \"Andreas\" <sip:andreas@e-horizon.se>;tag=8524\r\n"+
"Call-ID: f51a1851c5f570606140f14c8eb64fd3@172.18.1.29\r\n"+
"CSeq: 1 INVITE\r\n" +
"Max-Forwards: 70\r\n"+
"Record-Route: <sip:212.209.220.131:5060>\r\n"+
"Content-Length: 0\r\n\r\n",

"REGISTER sip:nist.gov SIP/2.0\r\n"+
"Via: SIP/2.0/UDP 129.6.55.182:14826\r\n"+
"Max-Forwards: 70\r\n"+
"From: <sip:mranga@nist.gov>;tag=6fcd5c7ace8b4a45acf0f0cd539b168b;epid=0d4c418ddf\r\n"+
"To: <sip:mranga@nist.gov>\r\n"+
"Call-ID: c5679907eb954a8da9f9dceb282d7230@129.6.55.182\r\n"+
"CSeq: 1 REGISTER\r\n"+
"Contact: <sip:129.6.55.182:14826>;methods=\"INVITE, MESSAGE, INFO, SUBSCRIBE, OPTIONS, BYE, CANCEL, NOTIFY, ACK, REFER\"\r\n"+
"User-Agent: RTC/(Microsoft RTC)\r\n"+
"Event:  registration\r\n"+
"Allow-Events: presence\r\n"+
"Content-Length: 0\r\n\r\n"+

    "INVITE sip:littleguy@there.com:5060 SIP/2.0\r\n"+
    "Via: SIP/2.0/UDP 65.243.118.100:5050\r\n" +
    "From: M. Ranganathan  <sip:M.Ranganathan@sipbakeoff.com>;tag=1234\r\n"+
    "To: \"littleguy@there.com\" <sip:littleguy@there.com:5060> \r\n" +
    "Call-ID: Q2AboBsaGn9!?x6@sipbakeoff.com \r\n" +
    "CSeq: 1 INVITE \r\n" +
    "Content-Length: 247\r\n\r\n"+
    "v=0\r\n"+
    "o=4855 13760799956958020 13760799956958020 IN IP4  129.6.55.78\r\n"+
    "s=mysession session\r\n"+
    "p=+46 8 52018010\r\n"+
    "c=IN IP4  129.6.55.78\r\n"+
    "t=0 0\r\n"+
    "m=audio 6022 RTP/AVP 0 4 18\r\n"+
    "a=rtpmap:0 PCMU/8000\r\n"+
    "a=rtpmap:4 G723/8000\r\n"+
    "a=rtpmap:18 G729A/8000\r\n"+
    "a=ptime:20\r\n"

        };

        for (int i = 0; i < messages.length; i++) {
            StringMsgParser smp = new StringMsgParser();
            SIPMessage sipMessage = smp.parseSIPMessage(messages[i]);
            System.out.println("encoded " + sipMessage.toString());
	    System.out.println("dialog id = " + sipMessage.GetDialogId(false));
        }


    }*/
