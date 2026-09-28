import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/library/search.mimetype.dropdown.dart';
import 'package:retrovibed/mimex.dart' as mimex;

void main() {
  group('SearchMimetypeDropdown', () {
    test('movies', () {
      final checksum = mimex.checksumfor(mimex.icomovie);
      expect(SearchMimetypeDropdown.mimetypesFor(checksum), mimex.of(mimex.icomovie));
      expect(SearchMimetypeDropdown.label(checksum), 'Movies');
      expect(SearchMimetypeDropdown.icon(checksum).icon, Icons.movie_filter);
    });

    // there is no unfiltered library; anything unrecognized is music.
    test('unrecognized falls back to music', () {
      final checksum = mimex.checksum(const []);
      expect(SearchMimetypeDropdown.mimetypesFor(checksum), mimex.of(mimex.icoaudio));
      expect(SearchMimetypeDropdown.label(checksum), 'Music');
      expect(SearchMimetypeDropdown.icon(checksum).icon, Icons.music_note);
    });
  });
}
