import 'dart:async';

import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/ddisc/api.dart' as ddisc;
import 'package:retrovibed/debug/diagnostics.details.dart';
import 'package:retrovibed/dhtx/api.dart' as dhtx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:retrovibed/torrentx/api.dart' as torrentx;

final _resolutions = Resolutions.variant();

Future<torrentx.TorrentMetricsResponse> _torrentNoop({List<httpx.Option> options = const []}) =>
    Completer<torrentx.TorrentMetricsResponse>().future;

Future<torrentx.TorrentMetricsResponse> _torrentEmpty({List<httpx.Option> options = const []}) async =>
    torrentx.TorrentMetricsResponse();

Future<torrentx.TorrentMetricsResponse> _torrentWithData({List<httpx.Option> options = const []}) async {
  return torrentx.TorrentMetricsResponse(
    torrent: torrentx.TorrentDiagnostics(total: fixnum.Int64(3)),
  );
}

Future<torrentx.TorrentMetricsResponse> _torrent404({List<httpx.Option> options = const []}) =>
    Future.error(http.Response('not found', 404));

Future<dhtx.DHTMetricsResponse> _dhtEmpty({List<httpx.Option> options = const []}) async => dhtx.DHTMetricsResponse();

Future<ddisc.DiscoveryMetricsResponse> _discoveryEmpty({List<httpx.Option> options = const []}) async =>
    ddisc.DiscoveryMetricsResponse();

void main() {
  group('DiagnosticsDetails', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        DiagnosticsDetails(
          apitorrent: _torrentEmpty,
          apidht: _dhtEmpty,
          apidiscovery: _discoveryEmpty,
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('shows loading indicator while fetching', (tester) async {
      await tester.pumpApp(
        physicalSize: const Size(1280, 720),
        DiagnosticsDetails(
          apitorrent: _torrentNoop,
          apidht: _dhtEmpty,
          apidiscovery: _discoveryEmpty,
        ),
      );
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders diagnostics after loading', (tester) async {
      await tester.pumpApp(
        physicalSize: const Size(1280, 720),
        DiagnosticsDetails(
          apitorrent: _torrentWithData,
          apidht: _dhtEmpty,
          apidiscovery: _discoveryEmpty,
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('3'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows not found error on 404', (tester) async {
      await tester.pumpApp(
        physicalSize: const Size(1280, 720),
        DiagnosticsDetails(
          apitorrent: _torrent404,
          apidht: _dhtEmpty,
          apidiscovery: _discoveryEmpty,
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('not found'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping the test-error button surfaces an error', (
      tester,
    ) async {
      await tester.pumpApp(
        physicalSize: const Size(1280, 720),
        DiagnosticsDetails(
          apitorrent: _torrentEmpty,
          apidht: _dhtEmpty,
          apidiscovery: _discoveryEmpty,
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.bug_report));
      await tester.pumpAndSettle();

      expect(find.text('an unexpected problem has occurred'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
