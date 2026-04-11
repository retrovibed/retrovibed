import 'package:flutter_test/flutter_test.dart';
import 'package:language_code/language_code.dart';
import 'package:retrovibed/langcodex.dart';

void main() {
  group('normalizeCode', () {
    test('returns correct code for plain language code', () {
      expect(normalizeCode('en'), LanguageCodes.en);
    });

    test('replaces hyphen with underscore', () {
      expect(normalizeCode('en-US'), LanguageCodes.en_US);
    });

    test('accepts underscore-separated code as-is', () {
      expect(normalizeCode('en_US'), LanguageCodes.en_US);
    });

    test('returns und for unknown code', () {
      expect(normalizeCode('xx'), LanguageCodes.und);
    });

    test('returns und for empty string', () {
      expect(normalizeCode(''), LanguageCodes.und);
    });
  });

  group('match', () {
    test('matches exact language', () {
      expect(match(LanguageCodes.en, 'en'), isTrue);
    });

    test('matches language prefix against regional variant', () {
      // LanguageCodes.en_US englishName starts with "English", same as en
      expect(match(LanguageCodes.en_US, 'en'), isTrue);
    });

    test('matches with hyphen in code', () {
      expect(match(LanguageCodes.en_US, 'en-US'), isTrue);
    });

    // regional variant tests
    test('en_GB matches base en', () {
      // "English (United Kingdom)".startsWith("English")
      expect(match(LanguageCodes.en_GB, 'en'), isTrue);
    });

    test('en_GB matches en_GB exactly', () {
      expect(match(LanguageCodes.en_GB, 'en_GB'), isTrue);
    });

    test('en_GB matches en-GB with hyphen', () {
      expect(match(LanguageCodes.en_GB, 'en-GB'), isTrue);
    });

    test('en_AU matches base en', () {
      // "English (Australia)".startsWith("English")
      expect(match(LanguageCodes.en_AU, 'en'), isTrue);
    });

    test('base en does not match en_GB', () {
      // "English".startsWith("English (United Kingdom)") is false
      expect(match(LanguageCodes.en, 'en_GB'), isFalse);
    });

    test('en_US does not match en_GB', () {
      // "English (United States of America)".startsWith("English (United Kingdom)") is false
      expect(match(LanguageCodes.en_US, 'en_GB'), isFalse);
    });

    test('en_GB does not match en_US', () {
      // "English (United Kingdom)".startsWith("English (United States of America)") is false
      expect(match(LanguageCodes.en_GB, 'en_US'), isFalse);
    });

    test('does not match different language', () {
      expect(match(LanguageCodes.fr, 'en'), isFalse);
    });

    test('does not match unknown code (und) against known language', () {
      expect(match(LanguageCodes.en, 'xx'), isFalse);
    });

    test('und matches und', () {
      expect(match(LanguageCodes.und, 'und'), isTrue);
    });
  });
}
