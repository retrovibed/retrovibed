import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:retrovibed/debug/metered.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/netmonx/api.dart' as api;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

Future<api.NetworkMetricsResponse> _noop({List<httpx.Option> options = const []}) =>
    Completer<api.NetworkMetricsResponse>().future;

Future<api.NetworkMetricsResponse> _empty({List<httpx.Option> options = const []}) async =>
    api.NetworkMetricsResponse();

Future<api.NetworkMetricsResponse> _withData({List<httpx.Option> options = const []}) async {
  return api.NetworkMetricsResponse(
    wireguard: api.WireguardDiagnostics(
      peerKey: 'peer-abc123',
      txBytes: ds.Int64(1500000000),
      rxBytes: ds.Int64(800000000),
      lastHandshakeSec: ds.Int64(DateTime.now().millisecondsSinceEpoch ~/ 1000),
      status: 'connected',
    ),
    network: api.Network(
      haveV4: true,
      haveV6: false,
      defaultInterface: 'eth0',
      interfaces: [
        api.NetworkInterface(name: 'eth0', ip: '192.168.1.10', metered: false),
        api.NetworkInterface(name: 'wlan0', ip: '10.0.0.5', metered: true),
      ],
    ),
  );
}

Future<api.NetworkMetricsResponse> _404({List<httpx.Option> options = const []}) =>
    Future.error(http.Response('not found', 404));

Future<api.NetworkMetricsResponse> _unauthorized({List<httpx.Option> options = const []}) =>
    Future.error(http.Response('unauthorized', 401));

void main() {
  group('MeteredToggle', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        const MeteredToggle(),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });

  group('MeteredCard', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        MeteredCard(apinetwork: _empty),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);

    testWidgets('shows loading indicator while fetching', (tester) async {
      await tester.pumpApp(
        physicalSize: const Size(1280, 720),
        MeteredCard(apinetwork: _noop),
      );
      await tester.pump();

      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders diagnostics after loading', (tester) async {
      await tester.pumpApp(
        physicalSize: const Size(1280, 720),
        MeteredCard(apinetwork: _withData),
      );
      await tester.pumpAndSettle();

      expect(find.text('connected'), findsOneWidget);
      expect(find.text('peer-abc123'), findsOneWidget);
      expect(find.text('eth0'), findsAtLeastNWidgets(1));
      expect(find.textContaining('192.168.1.10'), findsOneWidget);
      expect(find.textContaining('metered'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders placeholders when diagnostics are empty', (tester) async {
      await tester.pumpApp(
        physicalSize: const Size(1280, 720),
        MeteredCard(apinetwork: _empty),
      );
      await tester.pumpAndSettle();

      expect(find.text('never'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows not found error on 404', (tester) async {
      await tester.pumpApp(
        physicalSize: const Size(1280, 720),
        MeteredCard(apinetwork: _404),
      );
      await tester.pumpAndSettle();

      expect(find.text('not found'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows unauthorized error on 401', (tester) async {
      await tester.pumpApp(
        physicalSize: const Size(1280, 720),
        MeteredCard(apinetwork: _unauthorized),
      );
      await tester.pumpAndSettle();

      expect(find.text('you lack sufficient permissions'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
