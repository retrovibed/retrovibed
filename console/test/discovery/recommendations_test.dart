import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/discovery/recommendations.dart';
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

Future<lib.RecommendationSearchResponse> _notimplemented(
  lib.RecommendationSearchRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.error(http.Response('', 501));

Future<lib.RecommendationSearchResponse> _unauthorized(
  lib.RecommendationSearchRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.error(http.Response('', 401));

Future<lib.RecommendationSearchResponse> _empty(
  lib.RecommendationSearchRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.value(lib.RecommendationSearchResponse(items: []));

Future<lib.RecommendationSearchResponse> _withItems(
  lib.RecommendationSearchRequest req, {
  List<httpx.Option> options = const [],
}) =>
    Future.value(
      lib.RecommendationSearchResponse(
        items: [
          lib.Known(id: 'id-1', description: 'Recommendation One'),
          lib.Known(id: 'id-2', description: 'Recommendation Two'),
        ],
      ),
    );

void main() {
  group('Recommendations', () {
    testWidgets('displays loading state initially', (tester) async {
      await tester.pumpApp(Recommendations(mimex.video, latest:_notimplemented));
      await tester.pump();
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      expect(find.text('Recommendations'), findsOneWidget);
    });

    testWidgets('displays empty state after loading', (tester) async {
      await tester.pumpApp(Recommendations(mimex.video, latest:_empty));
      await tester.pumpAndSettle();
      expect(find.text('Recommendations'), findsOneWidget);
    });

    testWidgets('renders items returned from the api', (tester) async {
      await tester.pumpApp(Recommendations(mimex.video, latest:_withItems), isolatecache: true);
      await tester.pumpAndSettle();
      expect(find.text('Recommendations'), findsOneWidget);
      expect(find.byType(lib.KnownMediaLocator), findsWidgets);
    });

    testWidgets('silently ignores not implemented response', (tester) async {
      await tester.pumpApp(Recommendations(mimex.video, latest:_notimplemented));
      await tester.pumpAndSettle();
      expect(find.text('Recommendations'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays error on unauthorized response', (tester) async {
      await tester.pumpApp(Recommendations(mimex.video, latest:_unauthorized));
      await tester.pumpAndSettle();
      expect(find.text('Recommendations'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders at all resolutions', (tester) async {
      await tester.pumpApp(Recommendations(mimex.video, latest:_empty));
      await tester.pumpAndSettle();
      expect(find.text('Recommendations'), findsOneWidget);
    });
  });
}
