import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/library/grid.setting.dart';
import 'package:retrovibed/torrents.dart' as torrents;
import 'package:retrovibed/storage.dart' as storage;
import 'package:retrovibed/wireguard.dart' as wireguard;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('GridSettings', () {
    testWidgets('renders without overflow', (WidgetTester tester) async {
      final entry = _resolutions.currentValue!;
      await tester.pumpApp(
        GridSettings(
          torrents: ({options = const []}) => Future.value(torrents.TorrentSettings()),
          storage: ({options = const []}) => Future.value(storage.StorageSettingsResponse()),
          wgcurrent: () => Future.value(wireguard.WireguardCurrentResponse(wireguard: wireguard.Wireguard())),
          wgupdate:
              (v, {options = const []}) =>
                  Future.value(wireguard.WireguardUpdateResponse(wireguard: wireguard.Wireguard())),
        ),
        physicalSize: entry.value,
      );
      await tester.pumpAndSettle();
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });

  group('GridSettings inside ds.Grid', () {
    testWidgets('ds.Grid renders simple children', (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(
          width: 640,
          height: 560,
          child: ds.Grid<int>(
            (context, i) => Text('item_$i'),
            maxCrossAxisExtent: 320,
            aspectRatio: 1,
            children: [0, 1],
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('item_0'), findsOneWidget);
      expect(find.text('item_1'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('nested ds.Grid renders', (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(
          width: 640,
          height: 560,
          child: ds.Grid<Widget>(
            (context, child) => child,
            maxCrossAxisExtent: 320,
            aspectRatio: 4 / 3,
            children: [
              ds.Grid<int>(
                (context, i) => Container(
                  color: Colors.blue,
                  child: Center(child: Text('nested_$i')),
                ),
                maxCrossAxisExtent: 100,
                aspectRatio: 1,
                children: [1, 2],
              ),
            ],
          ),
        ),
      );

      await tester.pumpAndSettle();
      expect(find.text('nested_1'), findsOneWidget);
      expect(find.text('nested_2'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('GridSettings as ds.Grid child', (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(
          width: 500,
          height: 900,
          child: ds.Grid<Widget>(
            (context, child) => child,
            maxCrossAxisExtent: 600,
            aspectRatio: 7 / 8,
            children: [
              GridSettings(
                torrents: ({options = const []}) => Future.value(torrents.TorrentSettings()),
                storage: ({options = const []}) => Future.value(storage.StorageSettingsResponse()),
                wgcurrent: () => Future.value(wireguard.WireguardCurrentResponse(wireguard: wireguard.Wireguard())),
                wgupdate:
                    (v, {options = const []}) =>
                        Future.value(wireguard.WireguardUpdateResponse(wireguard: wireguard.Wireguard())),
              ),
            ],
          ),
        ),
      );

      await tester.pumpAndSettle();
      expect(find.byType(GridSettings), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('GridSettings with ds.Card', () {
    testWidgets('renders wrapped in ds.Card', (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(
          width: 500,
          height: 900,
          child: ds.Card(
            GridSettings(
              torrents: ({options = const []}) => Future.value(torrents.TorrentSettings()),
              storage: ({options = const []}) => Future.value(storage.StorageSettingsResponse()),
              wgcurrent: () => Future.value(wireguard.WireguardCurrentResponse(wireguard: wireguard.Wireguard())),
              wgupdate:
                  (v, {options = const []}) =>
                      Future.value(wireguard.WireguardUpdateResponse(wireguard: wireguard.Wireguard())),
            ),
          ),
        ),
      );

      await tester.pumpAndSettle();
      expect(find.byType(GridSettings), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
