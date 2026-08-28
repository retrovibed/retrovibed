import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/navbar/hamburger.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

List<String> _calls = [];
bool _maximized = false;

void _mockWindowManager() {
  _calls = [];
  _maximized = false;
  TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
      .setMockMethodCallHandler(
    const MethodChannel('window_manager'),
    (call) async {
      _calls.add(call.method);
      switch (call.method) {
        case 'isMaximized':
          return _maximized;
        case 'maximize':
          _maximized = true;
          return null;
        case 'unmaximize':
          _maximized = false;
          return null;
        case 'close':
          return null;
        default:
          return null;
      }
    },
  );
}

void _clearMock() {
  TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
      .setMockMethodCallHandler(
    const MethodChannel('window_manager'),
    null,
  );
}

void main() {
  setUp(_mockWindowManager);
  tearDown(_clearMock);

  Future<void> openMenu(WidgetTester tester) async {
    await tester.tap(find.byIcon(Icons.menu));
    await tester.pumpAndSettle();
  }

  group('Hamburger rendering', () {
    testWidgets('renders menu icon', (WidgetTester tester) async {
      await tester.pumpApp(Hamburger());
      await tester.pumpAndSettle();
      expect(find.byIcon(Icons.menu), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders as PopupMenuButton', (WidgetTester tester) async {
      await tester.pumpApp(Hamburger());
      await tester.pumpAndSettle();
      expect(find.byType(PopupMenuButton<String>), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Hamburger menu items', () {
    testWidgets('shows three items when opened', (WidgetTester tester) async {
      await tester.pumpApp(Hamburger());
      await tester.pumpAndSettle();
      await openMenu(tester);

      expect(find.text('Maximize'), findsOneWidget);
      expect(find.text('Logout'), findsOneWidget);
      expect(find.text('Exit'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows correct icons', (WidgetTester tester) async {
      await tester.pumpApp(Hamburger());
      await tester.pumpAndSettle();
      await openMenu(tester);

      expect(find.byIcon(Icons.fullscreen), findsOneWidget);
      expect(find.byIcon(Icons.logout), findsOneWidget);
      expect(find.byIcon(Icons.close), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows Minimize when maximized', (WidgetTester tester) async {
      _maximized = true;
      await tester.pumpApp(Hamburger());
      await tester.pumpAndSettle();
      await openMenu(tester);

      expect(find.text('Minimize'), findsOneWidget);
      expect(find.byIcon(Icons.fullscreen_exit), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Hamburger actions', () {
    testWidgets('maximize calls windowManager.maximize', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(Hamburger());
      await tester.pumpAndSettle();
      _calls.clear();
      await openMenu(tester);
      await tester.tap(find.text('Maximize'));
      await tester.pumpAndSettle();

      expect(_calls, contains('maximize'));
      expect(tester.takeException(), isNull);
    });

    testWidgets('minimize calls windowManager.unmaximize', (
      WidgetTester tester,
    ) async {
      _maximized = true;
      await tester.pumpApp(Hamburger());
      await tester.pumpAndSettle();
      _calls.clear();
      await openMenu(tester);
      await tester.tap(find.text('Minimize'));
      await tester.pumpAndSettle();

      expect(_calls, contains('unmaximize'));
      expect(tester.takeException(), isNull);
    });

    testWidgets('exit calls windowManager.close', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(Hamburger());
      await tester.pumpAndSettle();
      _calls.clear();
      await openMenu(tester);
      await tester.tap(find.text('Exit'));
      await tester.pumpAndSettle();

      expect(_calls, contains('close'));
      expect(tester.takeException(), isNull);
    });
  });
}
