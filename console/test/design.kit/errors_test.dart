import 'dart:async';
import 'dart:io';
import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

final _resolutions = Resolutions.variant();

void main() {
  group('ErrorTests.offline', () {
    test('detects connection refused (endpoint is down)', () {
      final obj = SocketException(
        'Connection refused',
        osError: OSError('Connection refused', 111),
      );
      expect(ds.ErrorTests.offline(obj), isTrue);
    });

    test('detects no route to host (endpoint does not exist)', () {
      final obj = SocketException(
        'No route to host',
        osError: OSError('No route to host', 113),
      );
      expect(ds.ErrorTests.offline(obj), isTrue);
    });

    test('detects invalid argument (daemon socket unreachable)', () {
      // matches the real-world crash: SocketException: Connection failed
      // (OS Error: Invalid argument, errno = 22), address = eg, port = 9998
      final obj = SocketException(
        'Connection failed',
        osError: OSError('Invalid argument', 22),
        address: InternetAddress('eg', type: InternetAddressType.unix),
        port: 9998,
      );
      expect(ds.ErrorTests.offline(obj), isTrue);
    });

    test('does not match unrelated os errors', () {
      final obj = SocketException(
        'Something else',
        osError: OSError('Something else', 1),
      );
      expect(ds.ErrorTests.offline(obj), isFalse);
    });

    test('does not match SocketException without osError', () {
      expect(ds.ErrorTests.offline(SocketException('no os error')), isFalse);
    });

    test('does not match non-SocketException errors', () {
      expect(ds.ErrorTests.offline(Exception('unrelated')), isFalse);
    });
  });

  group('ErrorTests.connectivity', () {
    test('detects HandshakeException', () {
      expect(ds.ErrorTests.connectivity(HandshakeException('tls failure')), isTrue);
    });

    test('does not match unrelated errors', () {
      expect(ds.ErrorTests.connectivity(Exception('unrelated')), isFalse);
    });
  });

  group('ErrorTests.dnsresolution', () {
    test('detects dns resolution failure', () {
      final obj = SocketException(
        'Failed host lookup',
        osError: OSError('nodename nor servname provided, or not known', -2),
      );
      expect(ds.ErrorTests.dnsresolution(obj), isTrue);
    });

    test('does not match connection refused', () {
      final obj = SocketException(
        'Connection refused',
        osError: OSError('Connection refused', 111),
      );
      expect(ds.ErrorTests.dnsresolution(obj), isFalse);
    });
  });

  group('ErrorTests.timeout', () {
    test('detects TimeoutException', () {
      expect(ds.ErrorTests.timeout(TimeoutException('timed out')), isTrue);
    });

    test('does not match unrelated errors', () {
      expect(ds.ErrorTests.timeout(Exception('unrelated')), isFalse);
    });
  });

  group('Error widget finite constraints', () {
    testWidgets('Error.zero renders as SizedBox in finite container', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Container(width: 200, height: 100, child: ds.Error.zero),
      );
      await tester.pumpAndSettle();

      expect(find.byType(SizedBox), findsWidgets);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Error widget renders in finite container', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Container(
          width: 200,
          height: 100,
          child: ds.Error.text('test error'),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('test error'), findsOneWidget);
      expect(find.byType(MouseRegion), findsWidgets);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Error widget expands to fill constrained parent', (
      WidgetTester tester,
    ) async {
      const width = 300.0;
      const height = 200.0;

      await tester.pumpApp(
        SizedBox(
          width: width,
          height: height,
          child: ds.Error.text('test error'),
        ),
      );
      await tester.pumpAndSettle();

      final size = tester.getSize(find.byType(ds.Error));
      expect(size.width, width);
      expect(size.height, height);
      expect(find.text('test error'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Error widget expands to fill parent across resolutions', (
      WidgetTester tester,
    ) async {
      final entry = _resolutions.currentValue!;
      final physicalSize = entry.value;

      await tester.pumpApp(
        ds.Error.text('test error'),
        physicalSize: physicalSize,
      );
      await tester.pumpAndSettle();

      final size = tester.getSize(find.byType(ds.Error));
      expect(size.width, physicalSize.width);
      expect(size.height, physicalSize.height);
      expect(tester.takeException(), isNull);
    }, variant: _resolutions);
  });

  group('Error widget infinite constraints', () {
    testWidgets('Error widget renders in ListView (infinite vertical)', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ListView(
          children: [
            ds.Error.text('error 1'),
            ds.Error.text('error 2'),
            ds.Error.zero,
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('error 1'), findsOneWidget);
      expect(find.text('error 2'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Error widget renders in Column (infinite vertical)', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: Column(
            children: [ds.Error.text('error in column'), ds.Error.zero],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('error in column'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Error widget renders in Row (infinite horizontal)', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: Row(
            children: [ds.Error.text('horizontal error'), ds.Error.zero],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('horizontal error'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Error widget factory methods', () {
    testWidgets('Error.text creates error with text', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error.text('factory text'),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('factory text'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Error.unknown creates error with default message', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error.unknown(Exception('test exception')),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('an unexpected problem has occurred'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Error.unauthorized creates error with permission message', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error.unauthorized(Exception('auth error')),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('you lack sufficient permissions'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Error.maybeErr returns Error.zero for null', (
      WidgetTester tester,
    ) async {
      final result = ds.Error.maybeErr(null);
      expect(result, equals(ds.Error.zero));
    });

    testWidgets('Error.maybeErr returns Error for exceptions', (
      WidgetTester tester,
    ) async {
      final result = ds.Error.maybeErr(Exception('test'));
      expect(result, isA<ds.Error>());
      expect(result, isNot(equals(ds.Error.zero)));
    });
  });

  group('Error widget with ErrorBoundary', () {
    testWidgets('ErrorBoundary displays error when onError called', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.ErrorBoundary(
          Container(width: 200, height: 100, child: Text('normal content')),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('normal content'), findsOneWidget);

      final boundary = tester.state<ds.ErrorBoundaryState>(
        find.byType(ds.ErrorBoundary),
      );
      boundary.onError(ds.Error.text('boundary error'));

      await tester.pumpAndSettle();

      expect(find.text('boundary error'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Error widget onTap', () {
    testWidgets('tapping error text triggers onTap', (
      WidgetTester tester,
    ) async {
      var tapped = false;

      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error.text('tap me', onTap: () => tapped = true),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('tap me'));
      await tester.pumpAndSettle();

      expect(tapped, isTrue);
    });

    testWidgets('error text is selectable', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error.text('select me', onTap: () {}),
        ),
      );
      await tester.pumpAndSettle();

      final gesture = await tester.createGesture(
        kind: PointerDeviceKind.mouse,
        pointer: 1,
      );
      await gesture.addPointer(
        location: tester.getCenter(find.text('select me')),
      );
      await tester.pump();

      expect(
        RendererBinding.instance.mouseTracker.debugDeviceActiveCursor(1),
        SystemMouseCursors.text,
      );
    });
  });

  group('Error widget long press', () {
    testWidgets('long press on error with cause shows dialog', (
      WidgetTester tester,
    ) async {
      final cause = Exception('something went wrong');

      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error.unknown(cause),
        ),
      );
      await tester.pumpAndSettle();

      await tester.longPress(find.text('an unexpected problem has occurred'));
      await tester.pumpAndSettle();

      expect(find.text('Error Details'), findsOneWidget);
      expect(find.text(cause.toString()), findsOneWidget);
      expect(find.text('Close'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('long press dialog can be dismissed', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error.unknown(Exception('dismiss me')),
        ),
      );
      await tester.pumpAndSettle();

      await tester.longPress(find.text('an unexpected problem has occurred'));
      await tester.pumpAndSettle();

      expect(find.text('Error Details'), findsOneWidget);

      await tester.tap(find.text('Close'));
      await tester.pumpAndSettle();

      expect(find.text('Error Details'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('dialog copy button copies diagnostic text to clipboard', (
      WidgetTester tester,
    ) async {
      final cause = Exception('copy me');
      String? copiedText;

      tester.binding.defaultBinaryMessenger.setMockMethodCallHandler(SystemChannels.platform, (
        call,
      ) async {
        if (call.method == 'Clipboard.setData') {
          copiedText = (call.arguments as Map)['text'] as String?;
        }
        return null;
      });
      addTearDown(
        () => tester.binding.defaultBinaryMessenger.setMockMethodCallHandler(
          SystemChannels.platform,
          null,
        ),
      );

      await tester.pumpApp(
        Scaffold(
          body: SizedBox(
            width: 200,
            height: 100,
            child: ds.Error.unknown(cause),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.longPress(find.text('an unexpected problem has occurred'));
      await tester.pumpAndSettle();

      await tester.tap(find.text('Copy'));
      await tester.pumpAndSettle();

      expect(copiedText, contains(cause.toString()));
      expect(find.text('Error details copied'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('long press on error without cause does nothing', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error.text('no cause here'),
        ),
      );
      await tester.pumpAndSettle();

      await tester.longPress(find.text('no cause here'));
      await tester.pumpAndSettle();

      expect(find.text('Error Details'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('long press shows full cause string in dialog', (
      WidgetTester tester,
    ) async {
      final cause = Exception('detailed error message');

      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error(
            cause: cause,
            child: const Text('user friendly message'),
            trace: StackTrace.current,
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.longPress(find.text('user friendly message'));
      await tester.pumpAndSettle();

      expect(find.text(cause.toString()), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Error widget rendering details', () {
    testWidgets('Error widget with custom color', (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error.text('colored error', color: Colors.blue),
        ),
      );
      await tester.pumpAndSettle();

      final container = tester.widget<Container>(
        find
            .descendant(
              of: find.byType(MouseRegion),
              matching: find.byType(Container),
            )
            .first,
      );

      final decoration = container.decoration as BoxDecoration;
      expect(decoration.color, Colors.blue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Error widget child is selectable', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: ds.Error.text('selectable error'),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(SelectionArea), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
