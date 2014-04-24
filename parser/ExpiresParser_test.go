package parser

import (
	"testing"
)

func TestExpiresParser(t *testing.T) {
	var tvs = []string{
		"Expires: 1000\n",
	}

	for i := 0; i < len(tvs); i++ {
		shp := NewExpiresParser(tvs[i])
		testHeaderParser(t, shp, tvs[i])
	}
}

/** Test program -- to be removed in final version.
    public static void main(String args[]) throws ParseException {
        String expires[] = {
            "Expires: 1000\n" };

            for (int i = 0; i < expires.length; i++ ) {
		try {
                	System.out.println("Parsing " + expires[i]);
                	ExpiresParser ep = new ExpiresParser(expires[i]);
                	Expires e = (Expires) ep.parse();
                	System.out.println("encoded = " +e.encode());
		} catch (ParseException ex) {
		  	System.out.println(ex.getMessage());
		}
            }

    }
*/
