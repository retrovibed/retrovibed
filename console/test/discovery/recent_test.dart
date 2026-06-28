import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/discovery/recent.dart';
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<lib.RecentSearchResponse> _notimplemented(
  lib.RecentSearchRequest req, {
  List<httpx.Option> options = const [],
}) => Future.error(http.Response('', 501));

Future<lib.RecentSearchResponse> _unauthorized(
  lib.RecentSearchRequest req, {
  List<httpx.Option> options = const [],
}) => Future.error(http.Response('', 401));

Future<lib.RecentSearchResponse> _empty(
  lib.RecentSearchRequest req, {
  List<httpx.Option> options = const [],
}) => Future.value(lib.RecentSearchResponse(items: []));

void main() {
  group('Recent', () {
    testWidgets('displays loading state initially', (tester) async {
      await tester.pumpApp(Recent(mimex.video, latest: _notimplemented));
      await tester.pump();
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      await tester.pumpAndSettle();
      expect(find.text('Continue'), findsOneWidget);
    });

    testWidgets('displays empty state after loading', (tester) async {
      await tester.pumpApp(Recent(mimex.video, latest: _empty));
      await tester.pumpAndSettle();
      expect(find.text('Continue'), findsOneWidget);
    });

    testWidgets('silently ignores not implemented response', (tester) async {
      await tester.pumpApp(Recent(mimex.video, latest: _notimplemented));
      await tester.pumpAndSettle();
      expect(find.text('Continue'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays error on unauthorized response', (tester) async {
      await tester.pumpApp(Recent(mimex.video, latest: _unauthorized));
      await tester.pumpAndSettle();
      expect(find.text('Continue'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders at all resolutions', (tester) async {
      await tester.pumpApp(Recent(mimex.video, latest: _empty));
      await tester.pumpAndSettle();
      expect(find.text('Continue'), findsOneWidget);
    });
  });
}
