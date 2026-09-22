import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/media/preview.text.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

http.StreamedResponse _streamed(String body) {
  return http.StreamedResponse(
    Stream.value(utf8.encode(body)),
    200,
  );
}

Future<http.StreamedResponse> _noop(String id, {List<httpx.Option> options = const []}) =>
    Completer<http.StreamedResponse>().future;

Future<http.StreamedResponse> _withContent(String id, {List<httpx.Option> options = const []}) async =>
    _streamed('line one\nline two\n');

Future<http.StreamedResponse> _truncated(String id, {List<httpx.Option> options = const []}) async =>
    _streamed('x' * PreviewText.limit);

Future<http.StreamedResponse> _unauthorized(String id, {List<httpx.Option> options = const []}) =>
    Future.error(http.Response('unauthorized', 401));

void main() {
  group('PreviewText', () {
    testWidgets('shows loading indicator while fetching', (tester) async {
      await tester.pumpApp(
        PreviewText(current: media.Media(id: 'm1'), download: _noop),
      );
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders the fetched content once loaded', (tester) async {
      await tester.pumpApp(
        PreviewText(current: media.Media(id: 'm1'), download: _withContent),
      );
      await tester.pumpAndSettle();

      expect(find.textContaining('line one'), findsOneWidget);
      expect(find.byType(PreviewText), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows the truncated notice when the body fills the limit', (
      tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: PreviewText(current: media.Media(id: 'm1'), download: _truncated),
        ),
      );
      await tester.pumpAndSettle();

      expect(
        find.text('preview truncated, download to read the rest'),
        findsOneWidget,
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows an error on unauthorized response', (tester) async {
      await tester.pumpApp(
        PreviewText(current: media.Media(id: 'm1'), download: _unauthorized),
      );
      await tester.pumpAndSettle();

      expect(find.text('you lack sufficient permissions'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
