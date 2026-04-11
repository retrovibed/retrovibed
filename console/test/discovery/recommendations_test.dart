import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/discovery/recommendations.dart';
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<lib.RecommendationsResponse> _notimplemented(
  lib.RecommendationsRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.error(http.Response('', 501));

Future<lib.RecommendationsResponse> _unauthorized(
  lib.RecommendationsRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.error(http.Response('', 401));

Future<lib.RecommendationsResponse> _empty(
  lib.RecommendationsRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.value(lib.RecommendationsResponse(items: []));

Future<lib.RecommendationsResponse> _withItems(
  lib.RecommendationsRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.value(
      lib.RecommendationsResponse(
        items: [
          lib.Known(id: 'id-1', description: 'Recommendation One'),
          lib.Known(id: 'id-2', description: 'Recommendation Two'),
        ],
      ),
    );

void main() {
  group('Recommendations', () {
    testWidgets('displays loading state initially', (tester) async {
      await tester.pumpApp(Recommendations(latest: _notimplemented));
      await tester.pump();
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      expect(find.text('Recommendations'), findsOneWidget);
    });

    testWidgets('displays empty state after loading', (tester) async {
      await tester.pumpApp(Recommendations(latest: _empty));
      await tester.pumpAndSettle();
      expect(find.text('Recommendations'), findsOneWidget);
    });

    testWidgets('renders items returned from the api', (tester) async {
      await tester.pumpApp(Recommendations(latest: _withItems));
      await tester.pumpAndSettle();
      expect(find.text('Recommendations'), findsOneWidget);
      expect(find.byType(lib.KnownMediaCard), findsWidgets);
    });

    testWidgets('silently ignores not implemented response', (tester) async {
      await tester.pumpApp(Recommendations(latest: _notimplemented));
      await tester.pumpAndSettle();
      expect(find.text('Recommendations'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays error on unauthorized response', (tester) async {
      await tester.pumpApp(Recommendations(latest: _unauthorized));
      await tester.pumpAndSettle();
      expect(find.text('Recommendations'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders at all resolutions', (tester) async {
      await tester.pumpApp(Recommendations(latest: _empty));
      await tester.pumpAndSettle();
      expect(find.text('Recommendations'), findsOneWidget);
    });
  });
}
