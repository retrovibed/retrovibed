import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/discovery/releases.dart';
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<lib.KnownLatestResponse> _notimplemented(
  lib.KnownLatestRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.error(http.Response('', 501));

Future<lib.KnownLatestResponse> _unauthorized(
  lib.KnownLatestRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.error(http.Response('', 401));

Future<lib.KnownLatestResponse> _empty(
  lib.KnownLatestRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.value(lib.KnownLatestResponse(items: []));

void main() {
  group('Releases', () {
    testWidgets('displays loading state initially', (tester) async {
      await tester.pumpApp(NewReleases(latest: _notimplemented));
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
    });

    testWidgets('displays empty state after loading', (tester) async {
      await tester.pumpApp(NewReleases(latest: _empty));
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
    });

    testWidgets('silently ignores not implemented response', (tester) async {
      await tester.pumpApp(NewReleases(latest: _notimplemented));
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays error on unauthorized response', (tester) async {
      await tester.pumpApp(NewReleases(latest: _unauthorized));
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders at all resolutions', (tester) async {
      await tester.pumpApp(NewReleases(latest: _empty));
      await tester.pumpAndSettle();
      expect(find.text('New Releases'), findsOneWidget);
    });
  });
}
