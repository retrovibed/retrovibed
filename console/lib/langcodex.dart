import 'package:language_code/language_code.dart';

LanguageCodes normalizeCode(String code) {
  return LanguageCodes.fromCode(
    code.replaceAll(
      "-",
      "_",
    ), // normalize the symbols used in langauges.
    orElse: () => LanguageCodes.und,
  );
}

bool match(LanguageCodes current, String code) {
  // return current.locale.languageCode == normalizeCode(code).locale.languageCode;
  return current.englishName.startsWith(normalizeCode(code).englishName);
}
