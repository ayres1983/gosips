package core

import (
    //"errors"
    "bytes"
)

/**
 * The base class from which all the other classes in the
 * sipheader, sdpfields and sipmessage packages are extended.
 * Provides a few utility funcitons such as indentation and
 * pretty printing that all other classes benifit from.
 *
 *@version  JAIN-SIP-1.1
 *
 *@author M. Ranganathan <mranga@nist.gov>  <br/>
 *
 *<a href="{@docRoot}/uncopyright.html">This code is in the public domain.</a>
 *
 */

type GenericObject interface {
    Clone() interface{}
    String() string
}

type GenericObjectImpl struct { // implements Serializable {
    indentation          int
    stringRepresentation string
    matchExpression      Match // Pattern matcher.
}

func NewGenericObjectImpl() *GenericObjectImpl{
	return &GenericObjectImpl{};
}

/** Set the  pattern matcher. To match on the
 * field of a sip message, set the match expression in the match template
 * and invoke the match function. This useful because
 * SIP headers and parameters may appear in different orders and are not
 * necessarily in canonical form. This makes it hard to write a pattern
 * matcher that relies on regular expressions alone.
 * Thus we rely on the following  strategy i.e. To do pattern matching on
 * an incoming message, first parse it, and then construct a match template,
 * filling in the fields that you want to
 * match. The rules for matching are: A null object matches wild card -
 * that is a match template of null matches any parsed SIP object.
 * To match with any subfield, set the match template on a template object
 * of the same type and invoke the match interface.
 * Regular expressions matching implements the gov.nist.sip.Match interface
 * that can be done using the Jakarta regexp package for example.
 * package included herein. This can be used to implement the Match interface
 * <a href=http://www.apache.org> See the APACHE website for documents </a>
 *
 */
func (this *GenericObjectImpl) SetMatcher(matchExpression Match) {
    //if (matchExpression == null)
    //    throw new IllegalArgumentException("null arg!");
    this.matchExpression = matchExpression
}

/** Return the match expression.
 *@return the match expression that has previously been set.
 */
func (this *GenericObjectImpl) GetMatcher() Match {
    return this.matchExpression
}

/*
   public static Class getClassFromName( String className) {
       try {
           return Class.forName(className);
       } catch ( Exception ex) {
           InternalErrorHandler.handleException(ex);
           return null;
       }
   }

   public static boolean isMySubclass(Class other) {
       try {
           return GenericObject.class.isAssignableFrom(other);
       } catch ( Exception ex) {
           InternalErrorHandler.handleException(ex);
       }
       return false;
   }
*/

/**
 *Make a clone of the given object.
 */
/*
    protected static Object makeClone( Object obj) {
	if (obj == null) throw new NullPointerException("null obj!");
        Object clone_obj = obj;
        if (obj instanceof String ) {
            String string = (String) obj;
            clone_obj = (Object) new String(string);
        } else if (obj instanceof Integer) {
            clone_obj = new Integer(((Integer)obj).intValue());
        } else if ( obj instanceof Float ) {
            clone_obj = new Float(((Float) obj).floatValue());
        } else if (obj instanceof Double) {
            clone_obj = new Double(((Double) obj).doubleValue());
        } else if (obj instanceof Long) {
            clone_obj = new Long(((Long)obj).longValue());
        } else {
            // If a clone method exists for the object, then
            // invoke it
            try {
                Class cl = obj.getClass();
                Method meth = cl.getMethod("clone",null);
                clone_obj = meth.invoke(obj,null);
            } catch (SecurityException ex) {
                clone_obj = obj;
            } catch (IllegalArgumentException ex) {
                InternalErrorHandler.handleException(ex);
            } catch (IllegalAccessException ex) {
                clone_obj = obj;
            } catch (InvocationTargetException ex) {
                clone_obj = obj;
            } catch (NoSuchMethodException ex) {
                clone_obj = obj;
            }
        }
        return clone_obj;
    }*/

/**
 *Make a clone of this object.
 */

func (this *GenericObjectImpl) Clone() interface{} {
    retval := &GenericObjectImpl{}
    retval.indentation = this.indentation
    retval.stringRepresentation = this.stringRepresentation
    retval.matchExpression = this.matchExpression // Pattern matcher.
    return retval
}

/**
 * Recursively override the fields of this object with the fields
 * of a new object. This is useful when you want to genrate a template
 * and override the fields of an incoming SIPMessage with another
 * SIP message that you have already generated.
 *
 * @param mergeObject is the replacement object.  The override
 * obect must be of the same class as this object.
 * Set any fields that you do not want to override as null in the
 * mergeOject object.
 */
/*
    public void merge(Object mergeObject) {
        if (! mergeObject.getClass().equals(this.getClass()))
            throw new IllegalArgumentException("Bad override object");
        // Base case.
        if (mergeObject == null) return;
        Class myclass = this.getClass();
	while (true) {
        Field[] fields = myclass.getDeclaredFields();
        for (int i = 0; i < fields.length; i++) {
            Field f = fields [i];
            int modifier =  f.getModifiers();
            if (Modifier.isPrivate(modifier)) {
                continue;
            } else if (Modifier.isStatic(modifier)) {
                continue;
            } else if (Modifier.isInterface(modifier)) {
                continue;
            }
            Class fieldType = f.getType();
            String fieldName = f.getName();
            String fname = fieldType.toString();
            try {
                // Primitive fields are printed with type: value
                if (fieldType.isPrimitive()) {
                    if (fname.compareTo("int") == 0) {
                        int intfield = f.getInt(mergeObject);
                        f.setInt(this,intfield);
                    } else if (fname.compareTo("short") == 0) {
                        short shortField = f.getShort(mergeObject);
                        f.setShort(this,shortField);
                    } else if (fname.compareTo("char") == 0) {
                        char charField = f.getChar(mergeObject);
                        f.setChar(this,charField);
                    } else if (fname.compareTo("long") == 0) {
                        long longField = f.getLong(mergeObject);
                        f.setLong(this,longField);
                    } else if (fname.compareTo("boolean") == 0) {
                        boolean booleanField = f.getBoolean(mergeObject);
                        f.setBoolean(this,booleanField);
                    } else if (fname.compareTo("double") == 0) {
                        double doubleField = f.getDouble(mergeObject);
                        f.setDouble(this,doubleField);
                    } else if (fname.compareTo("float") == 0) {
                        float floatField = f.getFloat(mergeObject);
                        f.setFloat(this,floatField);
                    }
                } else {
                    Object obj = f.get(this);
                    Object mobj = f.get(mergeObject);
                    if (mobj == null) continue;
                    if (obj == null) {
                        f.set(this,mobj);
                        continue;
                    }
                    if (obj instanceof GenericObject)  {
                        GenericObject gobj = (GenericObject) obj;
                        gobj.merge(mobj);
                    } else {
                        f.set(this,mobj);
                    }
                }
            }  catch (IllegalAccessException ex1 ) {
                ex1.printStackTrace();
                continue; // we are accessing a private field...
            }
        }
	myclass = myclass.getSuperclass();
	if (myclass.equals(GenericObject.class)) break;
	}
    }
*/

/*protected GenericObject() {
    indentation =  0;
    stringRepresentation = "";
}*/

func (this *GenericObjectImpl) GetIndentation() string {
    var retval bytes.Buffer
    for i := 0; i < this.indentation; i++ {
        retval.WriteString(" ");
    }
    return retval.String();
}

/**
 * Add a new string to the accumulated string representation.
 */

func (this *GenericObjectImpl) Sprint(a string) {
    if a == "" {
        this.stringRepresentation += this.GetIndentation()
        this.stringRepresentation += "<null>\n"
        return
    }
    if a == "}" || a == "]" {
        this.indentation--
    }
    this.stringRepresentation += this.GetIndentation()
    this.stringRepresentation += a
    this.stringRepresentation += "\n"
    if a == "{" || a == "[" {
        this.indentation++
    }
}

/**
 * Pretty printing function accumulator for objects.
 */

/*protected void sprint( Object o ) {
    sprint(o.toString());
}*/

/**
 * Pretty printing accumulator function for ints
 */

/*protected void sprint( int intField) {
    sprint((new Integer(intField)).toString());
}*/

/**
 * Pretty printing accumulator function for shorts
 */
/*protected void sprint( short shortField) {
    sprint((new Short(shortField)).toString());
}*/

/**
 * Pretty printing accumulator function for chars
 */

/*protected void sprint( char charField) {
    sprint((new Character(charField)).toString());
}*/

/**
 * Pretty printing accumulator function for longs
 */

/*protected void sprint(long longField) {
    sprint((new Long(longField)).toString());
}*/

/**
 * Pretty printing accumulator function for booleans
 */

/*protected void sprint(boolean booleanField) {
    sprint((new Boolean(booleanField)).toString());
}*/

/**
 * Pretty printing accumulator function for doubles
 */

/*protected void sprint(double doubleField) {
    sprint((new Double(doubleField)).toString());
}*/

/**
 * Pretty printing accumulator function for floats
 */

/*protected void sprint( float floatField) {
    sprint((new Float(floatField)).toString());
}*/

/**
 * Debug printing function.
 */

func (this *GenericObjectImpl) dbgPrint() {
    Debug.println(this.debugDump())
}

/**
 * Debug printing function.
 */
func (this *GenericObjectImpl) dbgPrintln(s string) {
    Debug.println(s)
}

/**
 * An introspection based equality predicate for GenericObjects.
 *@param other is the other object to test against.
 *@return true if the objects are euqal and false otherwise
 */
/*public boolean equals(Object that) {
        if (! this.getClass().equals(that.getClass())) return false;
        Class myclass = this.getClass();
        Class hisclass = that.getClass();
	while (true)  {
        Field[] fields = myclass.getDeclaredFields();
        Field[] hisfields = hisclass.getDeclaredFields();
        for (int i = 0; i < fields.length; i++) {
            Field f = fields [i];
            Field g = hisfields[i];
            // Only print protected and public members.
            int modifier =  f.getModifiers();
            if ( (modifier & Modifier.PRIVATE) == Modifier.PRIVATE)
                continue;
            Class fieldType = f.getType();
            String fieldName = f.getName();
            if (fieldName.compareTo("stringRepresentation") == 0 ) {
                continue;
            }
            if (fieldName.compareTo("indentation") == 0 ) {
                continue;
            }
            try {
                // Primitive fields are printed with type: value
                if (fieldType.isPrimitive()) {
                    String fname = fieldType.toString();
                    if (fname.compareTo("int") == 0) {
                        if (f.getInt(this) != g.getInt(that)) return false;
                    } else if (fname.compareTo("short") == 0) {
                        if (f.getShort(this) != g.getShort(that)) return false;
                    } else if (fname.compareTo("char") == 0) {
                        if (f.getChar(this) != g.getChar(that)) return false;
                    } else if (fname.compareTo("long") == 0) {
                        if (f.getLong(this) != g.getLong(that)) return false;
                    } else if (fname.compareTo("boolean") == 0) {
                        if (f.getBoolean(this) != g.getBoolean(that))
                            return false;
                    } else if (fname.compareTo("double") == 0) {
                        if (f.getDouble(this) != g.getDouble(that))
                            return false;
                    } else if (fname.compareTo("float") == 0) {
                        if (f.getFloat(this) != g.getFloat(that)) return false;
                    }
                } else if (g.get(that) == f.get(this)) return true;
                else if (f.get(this) == null) return false;
                else if (g.get(that) == null) return false;
                else if (g.get(that) == null && f.get(this) != null )
                        return false;
                else if (! f.get(this).equals(g.get(that))) return false;
            }  catch (IllegalAccessException ex1 ) {
                InternalErrorHandler.handleException(ex1);
            }
        }
        if (myclass.equals(GenericObject.class)) break;
        else {
		myclass = myclass.getSuperclass();
		hisclass = hisclass.getSuperclass();
	}

	}
        return true;
    }*/

/** An introspection based predicate matching using a template
 * object. Allows for partial match of two protocl Objects.
 *@param other the match pattern to test against. The match object
 * has to be of the same type (class). Primitive types
 * and non-sip fields that are non null are matched for equality.
 * Null in any field  matches anything. Some book-keeping fields
 * are ignored when making the comparison.
 */

/*public boolean match(Object other) {
    if (other == null) return true;
    if ( ! this.getClass().equals(other.getClass()))
        return false;
    GenericObject that = (GenericObject)other;
    Class myclass = this.getClass();
    Field[] fields = myclass.getDeclaredFields();
    Class hisclass = other.getClass();
    Field[] hisfields = hisclass.getDeclaredFields();
    for (int i = 0; i < fields.length; i++) {
        Field f = fields [i];
        Field g = hisfields[i];
        // Only print protected and public members.
        int modifier =  f.getModifiers();
        if ( (modifier & Modifier.PRIVATE) == Modifier.PRIVATE)
            continue;
        Class fieldType = f.getType();
        String fieldName = f.getName();
        if (fieldName.compareTo("stringRepresentation") == 0 ) {
            continue;
        }
        if (fieldName.compareTo("indentation") == 0 ) {
            continue;
        }
        try {
            // Primitive fields are printed with type: value
            if (fieldType.isPrimitive()) {
                String fname = fieldType.toString();
                if (fname.compareTo("int") == 0) {
                    if (f.getInt(this) != g.getInt(that)) return false;
                } else if (fname.compareTo("short") == 0) {
                    if (f.getShort(this) != g.getShort(that)) return false;
                } else if (fname.compareTo("char") == 0) {
                    if (f.getChar(this) != g.getChar(that)) return false;
                } else if (fname.compareTo("long") == 0) {
                    if (f.getLong(this) != g.getLong(that)) return false;
                } else if (fname.compareTo("boolean") == 0) {
                    if (f.getBoolean(this) != g.getBoolean(that))
                        return false;
                } else if (fname.compareTo("double") == 0) {
                    if (f.getDouble(this) != g.getDouble(that))
                        return false;
                } else if (fname.compareTo("float") == 0) {
                    if (f.getFloat(this) != g.getFloat(that)) return false;
                }
            } else {
                Object myObj = f.get(this);
                Object hisObj = g.get(that);
                if (hisObj != null && myObj == null) return false;
                else if (hisObj == null && myObj != null) continue;
                else if (hisObj == null && myObj == null) continue;
                else if (hisObj instanceof java.lang.String &&
                myObj instanceof java.lang.String) {
                    if ((((String)hisObj).trim()).equals("")) continue;
                    if (((String)myObj).compareToIgnoreCase
                    ((String)hisObj) != 0) return false;
                }
                else if (GenericObject.isMySubclass(myObj.getClass())
                && ! ((GenericObject) myObj).match(hisObj))
                    return false;
                else if (GenericObjectList.isMySubclass(myObj.getClass())
                && ! ((GenericObjectList) myObj).match(hisObj))
                    return false;

            }
        }  catch (IllegalAccessException ex1 ) {
            InternalErrorHandler.handleException(ex1);
        }
    }
    return true;
}*/

/**
 * Generic print formatting function:
 * Does depth-first descent of the structure and
 * recursively prints all non-private objects pointed to
 * by this object.
 * <bf>
 * Warning - the following generic string routine will
 * bomb (go into infinite loop) if there are any circularly linked
 * structures so if you have these, they had better be private!
 * </bf>
 * We dont have to worry about such things for our structures
 *(we never use circular linked structures).
 */

func (this *GenericObjectImpl) debugDump() string {
    /*
       stringRepresentation = "";
       Class myclass = getClass();
       sprint(myclass.getName());
       sprint("{");
       Field[] fields = myclass.getDeclaredFields();
       for (int i = 0; i < fields.length; i++) {
           Field f = fields [i];
           // Only print protected and public members.
           int modifier =  f.getModifiers();
           if ( (modifier & Modifier.PRIVATE)  == Modifier.PRIVATE)
               continue;
           Class fieldType = f.getType();
           String fieldName = f.getName();
           if (fieldName.compareTo("stringRepresentation") == 0 ) {
               // avoid nasty recursions...
               continue;
           }
           if (fieldName.compareTo("indentation") == 0 ) {
               // formatting stuff - not relevant here.
               continue;
           }
           sprint(fieldName + ":"  );
           try {
               // Primitive fields are printed with type: value
               if (fieldType.isPrimitive()) {
                   String fname = fieldType.toString();
                   sprint(fname + ":");
                   if (fname.compareTo("int") == 0) {
                       int intfield = f.getInt(this);
                       sprint(intfield);
                   } else if (fname.compareTo("short") == 0) {
                       short shortField = f.getShort(this);
                       sprint(shortField);
                   } else if (fname.compareTo("char") == 0) {
                       char charField = f.getChar(this);
                       sprint(charField);
                   } else if (fname.compareTo("long") == 0) {
                       long longField = f.getLong(this);
                       sprint(longField);
                   } else if (fname.compareTo("boolean") == 0) {
                       boolean booleanField = f.getBoolean(this);
                       sprint(booleanField);
                   } else if (fname.compareTo("double") == 0) {
                       double doubleField = f.getDouble(this);
                       sprint(doubleField);
                   } else if (fname.compareTo("float") == 0) {
                       float floatField = f.getFloat(this);
                       sprint(floatField);
                   }
               } else if (GenericObject.class.isAssignableFrom(fieldType)) {
                   if (f.get(this) != null) {
                       sprint(((GenericObject)f.get(this)).
                       debugDump(indentation + 1));
                   } else {
                       sprint("<null>");
                   }

               } else if
               ( GenericObjectList.class.isAssignableFrom(fieldType))  {
                   if (f.get(this) != null) {
                       sprint(((GenericObjectList)f.get(this)).
                       debugDump(indentation + 1));
                   } else {
                       sprint("<null>");
                   }

               }  else {
                   // Dont do recursion on things that are not
                   // of our header type...
                   if (f.get(this) != null ) {
                       sprint(f.get(this).getClass().getName() + ":");
                   } else {
                       sprint(fieldType.getName() + ":");
                   }

                   sprint("{");
                   if (f.get(this) != null) {
                       sprint(f.get(this).toString());
                   } else {
                       sprint("<null>");
                   }
                   sprint("}");
               }
           }  catch (IllegalAccessException ex1 ) {
               continue; // we are accessing a private field...
           } catch (Exception ex) {
               InternalErrorHandler.handleException(ex);
           }
       }
       sprint("}");
       return stringRepresentation;
    */
    return ""
}

/**
 * Formatter with a given starting indentation.
 */
func (this *GenericObjectImpl) DebugDump(indent int) string {
    this.indentation = indent
    retval := this.debugDump()
    this.indentation = 0
    return retval
}

/**
 * An assertion checking utility.
 */
/*
func (this *GenericObjectImpl) Assert(condition bool, msg string) error {
    if !condition {
        //InternalErrorHandler.handleException(msg);
        return errors.New(msg)
    }

    return nil
}*/

/**
 *  Get the string encoded version of this object
 * @since v1.0
 */
func (this *GenericObjectImpl) String() string {
    return ""
}

/**
 * Do a recursive find and replace of objects pointed to by this
 * object.
 * @since v1.0
 * @param objectText is the canonical string representation of
 *		the object that we want to replace.
 * @param replacement is the object that we want to replace it
 *		with.
 * @param matchSubstring a boolean which tells if we should match
 * 		a substring of the target object
 * A replacement will occur if a portion of the structure is found
 * with matching encoded text (a substring if matchSubstring is true)
 * as objectText and with the same class as replacement.
 */
/*public void replace(String objectText,
  GenericObject replacement,
  boolean matchSubstring )
  throws IllegalArgumentException {
      if (objectText == null || replacement == null) {
          throw new IllegalArgumentException("null argument!");
      }
      Class replacementClass = replacement.getClass();
      Class myclass = getClass();
      Field[] fields = myclass.getDeclaredFields();
      for (int i = 0; i < fields.length; i++) {
          Field f = fields[i];
          Class fieldType = f.getType();
          if (! GenericObject.class.isAssignableFrom(fieldType) &&
          ! GenericObjectList.class.isAssignableFrom(fieldType)) {
              continue;
          } else if ( (f.getModifiers() & Modifier.PRIVATE)
          == Modifier.PRIVATE) {
              continue;
          }

          try {
              if (fieldType.equals(replacementClass)) {

                  if (GenericObject.isMySubclass(replacementClass)) {
                      GenericObject  obj =
                      (GenericObject)f.get(this);
                      if (! matchSubstring) {
                          if (objectText.compareTo(obj.encode())== 0 ) {
                              f.set(this,replacement);
                          }
                      } else {
                          // Substring match is specified
                          if (obj.encode().indexOf(objectText) >= 0 ) {
                              f.set(this,replacement);
                          }
                      }
                  }
              } else if (GenericObjectList.isMySubclass
              (replacementClass)) {
                  GenericObjectList  obj =
                  (GenericObjectList)f.get(this);
                  if (! matchSubstring) {
                      if (objectText.compareTo(obj.encode())== 0 ) {
                          f.set(this,replacement);
                      }
                  } else  {
                      if (obj.encode().indexOf(objectText) >= 0 ) {
                          f.set(this,replacement);
                      }
                  }
              } else if ( GenericObject.class
              .isAssignableFrom(fieldType)){
                  GenericObject g =
                  (GenericObject) f.get(this);
                  g.replace(objectText, replacement,matchSubstring);
              } else if (GenericObjectList.class
              .isAssignableFrom(fieldType)){
                  GenericObjectList g =
                  (GenericObjectList) f.get(this);
                  g.replace( objectText, replacement,matchSubstring);
              }
          } catch (IllegalAccessException ex) {
              InternalErrorHandler.handleException(ex);
          }
      }

  }*/

/**
 * Do a recursive find and replace of objects pointed to by this
 * object.
 * @since v1.0
 *@param objectText Canonical string representation of the
 *  portion we want to replace.
 *@param replacement object we want to replace this portion with.
 * A replacement will occur if a portion of the structure is found
 * with a match of the  encoded text
 * with objectText and with the same class as replacement.
 *@param matchSubstring is true if we want to match objectText
 * 	as a substring of the encoded target text.
 * (i.e. an object is a  candidate for replacement if
 *   objectText is a substring of
 *  candidate.encode() && candidate.class.equals(replacement.class)
 *  otherwise the match test is an equality test.)
 */
/*public void replace(String objectText,
  GenericObjectList replacement,
  boolean matchSubstring )
  throws IllegalArgumentException {

      if (objectText == null || replacement == null) {
          throw new IllegalArgumentException("null argument!");
      }
      Class replacementClass = replacement.getClass();
      Class myclass = getClass();
      Field[] fields = myclass.getDeclaredFields();
      for (int i = 0; i < fields.length; i++) {
          Field f = fields[i];
          Class fieldType = f.getType();
          if (!  GenericObject.class.isAssignableFrom(fieldType) &&
          !  GenericObjectList.class.isAssignableFrom(fieldType)) {
              continue;
          } else if ( (f.getModifiers() & Modifier.PRIVATE)
          == Modifier.PRIVATE) {
              continue;
          }
          try {
              if (fieldType.equals(replacementClass)) {
                  if (GenericObject.isMySubclass(replacementClass)) {
                      GenericObject  obj =
                      (GenericObject)f.get(this);
                      if (! matchSubstring) {
                          if (objectText.compareTo(obj.encode())== 0 ) {
                              f.set(this,replacement);
                          }
                      } else {
                          if (obj.encode().indexOf(objectText) >= 0 ) {
                              f.set(this,replacement);
                          }
                      }
                  } else if (GenericObjectList.isMySubclass
                  (replacementClass)) {
                      GenericObjectList  obj =
                      (GenericObjectList)f.get(this);
                      if (! matchSubstring) {
                          if (objectText.compareTo(obj.encode())== 0 ) {
                              f.set(this,replacement);
                          }
                      } else  {
                          if (obj.encode().indexOf(objectText) >= 0 ) {
                              f.set(this,replacement);
                          }
                      }
                  }

              } else if ( GenericObject.class
              .isAssignableFrom(fieldType)){
                  GenericObject g = (GenericObject) f.get(this);
                  g.replace(objectText, replacement,matchSubstring);
              }else if (
              GenericObjectList.class
              .isAssignableFrom(fieldType)){
                  GenericObjectList g =
                  (GenericObjectList) f.get(this);
                  g.replace(objectText, replacement,matchSubstring);
              }
          } catch (IllegalAccessException ex) {
              InternalErrorHandler.handleException(ex);
          }
      }

  }*/
/**
 * Do a recursive find and replace of objects pointed to by this
 * object based on regular expression pattern matching.
 * @since v1.0
 *@param regexp  regular expression for the object we want to find.
 * This is generated using a regular expression matching package
 * such as the apache regexp package.
 *@param replacement object we want to replace this portion with.
 * A replacement will occur if a portion of the structure is found
 * with a match of the  encoded text
 * with objectText and with the same class as replacement.
 */
/*public void replace(Match regexp,
    GenericObjectList replacement)
    throws IllegalArgumentException {

        if (regexp == null || replacement == null) {
            throw new IllegalArgumentException("null argument!");
        }
        Class replacementClass = replacement.getClass();
        Class myclass = getClass();
        Field[] fields = myclass.getDeclaredFields();
        for (int i = 0; i < fields.length; i++) {
            Field f = fields[i];
            Class fieldType = f.getType();
            if (!  GenericObject.class.isAssignableFrom(fieldType) &&
            !  GenericObjectList.class.isAssignableFrom(fieldType)) {
                continue;
            } else if ( (f.getModifiers() & Modifier.PRIVATE)
            == Modifier.PRIVATE) {
                continue;
            }
            try {
                if (fieldType.equals(replacementClass)) {
                    if (GenericObject.isMySubclass(replacementClass)) {
                        GenericObject  obj =
                        (GenericObject)f.get(this);
                        if (regexp.match(obj.encode())) f.set(this,replacement);
                    } else if (GenericObjectList.isMySubclass
                    (replacementClass)) {
                        GenericObjectList  obj =
                        (GenericObjectList)f.get(this);
                        if (regexp.match(obj.encode())) f.set(this,replacement);
                    }

                } else if ( GenericObject.class.isAssignableFrom(fieldType)){
                    GenericObject g = (GenericObject) f.get(this);
                    g.replace(regexp, replacement);
                }else if ( GenericObjectList.class.isAssignableFrom(fieldType)){
		    GenericObjectList g =
                    (GenericObjectList) f.get(this);
                    g.replace(regexp, replacement);
                }
            } catch (IllegalAccessException ex) {
                InternalErrorHandler.handleException(ex);
            }
        }

    }*/

/**
 * Do a find and replace of objects based on regular expression
 * matching of fields.
 * @param regexp is the match expression (i.e. implementation of
 *		the Match interface) for
 *		the object that we want to replace.
 * @param replacement is the object that we want to replace it
 *		with.
 * A replacement will occur if a portion of the structure is found
 * that matches according to the given regexp and if the class of
 * the replaced field matches the replacement.
 */
/*public void replace(Match regexp, GenericObject replacement)
  throws IllegalArgumentException {
      if (regexp == null || replacement == null) {
          throw new IllegalArgumentException("null argument!");
      }
      Class replacementClass = replacement.getClass();
      Class myclass = getClass();
      Field[] fields = myclass.getDeclaredFields();
      for (int i = 0; i < fields.length; i++) {
          Field f = fields[i];
          Class fieldType = f.getType();
          if (! GenericObject.class.isAssignableFrom(fieldType) &&
              ! GenericObjectList.class.isAssignableFrom(fieldType)) {
              continue;
          } else if (f.getModifiers() == Modifier.PRIVATE) {
              continue;
          }

          try {
              if (fieldType.equals(replacementClass)) {

                  if (GenericObject.isMySubclass(replacementClass)) {
                      GenericObject  obj = (GenericObject)f.get(this);
                      if (regexp.match(obj.encode())) f.set(this,replacement);
                  }
              } else if (GenericObjectList.isMySubclass
              (replacementClass)) {
                  GenericObjectList  obj = (GenericObjectList)f.get(this);
                  if (regexp.match(obj.encode())) f.set(this,replacement);
              } else if (GenericObject.class.isAssignableFrom(fieldType)){
                  GenericObject g =
                  (GenericObject) f.get(this);
                  g.replace(regexp, replacement);
              } else if (GenericObjectList.class.isAssignableFrom(fieldType)){
                  GenericObjectList g =
                  (GenericObjectList) f.get(this);
                  g.replace( regexp, replacement);
              }
          } catch (IllegalAccessException ex) {
              InternalErrorHandler.handleException(ex);
          }
      }

  }*/