package text

import "strings"

func Filter(s string) string {
	output := strings.Replace(s, string(rune(1)), "", -1)
	output = strings.Replace(s, string(rune(2)), "", -1)
	output = strings.Replace(s, string(rune(9)), "", -1)
	output = strings.Replace(s, string(rune(10)), "", -1)
	output = strings.Replace(s, string(rune(12)), "", -1)
	output = strings.Replace(s, string(rune(13)), "", -1) // filter newline chars too
	return output
}

// 1234567890qwertyuiopasdfghjklzxcvbnm_-+=?!@:.,$
func ContainsAllowedChars(toTest, allowedChars string) bool {
	for _, v := range toTest {
		if strings.Contains(allowedChars, string(v)) {
			continue
		}
		return false
	}

	return true
}

/*
   public static boolean hasAllowedCharacters(String str, String allowedChars) {
       if (str == null) {
           return false;
       }

       for (int i = 0; i < str.length(); i++) {
           if (allowedChars.contains(Character.valueOf(str.toCharArray()[i]).toString())) {
               continue;
           }

           return false;
       }

       return true;
   }
*/
