import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:plugin_platform_interface/plugin_platform_interface.dart';
import 'package:retrovibed/design.kit/hyperlink.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';
import 'package:url_launcher_platform_interface/url_launcher_platform_interface.dart';

class _FakeUrlLauncher extends Fake
    with MockPlatformInterfaceMixin
    implements UrlLauncherPlatform {
  final List<String> launched = [];

  @override
  Future<bool> launchUrl(String url, LaunchOptions options) async {
    launched.add(url);
    return true;
  }
}

void main() {
  group('Hyperlink rendering', () {
    testWidgets('renders text via Uri constructor', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Hyperlink(
          'Click here',
          uri: Uri.parse('https://example.com'),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Click here'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders text via fromString constructor', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Hyperlink.fromString(
          'Click here',
          url: 'https://example.com',
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Click here'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('applies underline decoration and primary color', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Hyperlink.fromString(
          'styled link',
          url: 'https://example.com',
        ),
        theme: ThemeData(colorScheme: ColorScheme.light()),
      );
      await tester.pumpAndSettle();

      final text = tester.widget<Text>(find.text('styled link'));
      expect(text.style?.decoration, TextDecoration.underline);
      expect(
        text.style?.color,
        ThemeData(colorScheme: ColorScheme.light()).colorScheme.primary,
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('falls back to theme bodyMedium when no custom style provided',
        (WidgetTester tester) async {
      final theme = ThemeData(
        textTheme: TextTheme(
          bodyMedium: TextStyle(fontSize: 18, fontWeight: FontWeight.w300),
        ),
      );

      await tester.pumpApp(
        Hyperlink.fromString(
          'default style',
          url: 'https://example.com',
        ),
        theme: theme,
      );
      await tester.pumpAndSettle();

      final text = tester.widget<Text>(find.text('default style'));
      expect(text.style?.fontSize, 18);
      expect(text.style?.fontWeight, FontWeight.w300);
      expect(tester.takeException(), isNull);
    });

    testWidgets('uses custom style when provided', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Hyperlink.fromString(
          'custom style',
          url: 'https://example.com',
          style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
        ),
      );
      await tester.pumpAndSettle();

      final text = tester.widget<Text>(find.text('custom style'));
      expect(text.style?.fontSize, 24);
      expect(text.style?.fontWeight, FontWeight.bold);
      expect(text.style?.decoration, TextDecoration.underline);
      expect(tester.takeException(), isNull);
    });
  });

  group('Hyperlink interaction', () {
    testWidgets('tapping calls launchUrl with correct URI', (
      WidgetTester tester,
    ) async {
      final fake = _FakeUrlLauncher();
      UrlLauncherPlatform.instance = fake;

      await tester.pumpApp(
        Hyperlink.fromString(
          'tap me',
          url: 'https://example.com/page',
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('tap me'));
      await tester.pumpAndSettle();

      expect(fake.launched, contains('https://example.com/page'));
      expect(tester.takeException(), isNull);
    });
  });

  group('Hyperlink constrained parents', () {
    testWidgets('renders in a tight SizedBox', (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 40,
          child: Hyperlink.fromString(
            'short',
            url: 'https://example.com',
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('short'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('long text in Row without Flexible overflows', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 50,
          child: Row(
            children: [
              Hyperlink.fromString(
                'This is an extremely long hyperlink text that will certainly overflow',
                url: 'https://example.com',
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(Hyperlink), findsOneWidget);

      final exception = tester.takeException();
      expect(exception, isA<FlutterError>());
      expect(exception.toString(), contains('overflowed'));
    });

    testWidgets('long text in constrained Row with Flexible wraps without error',
        (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(
          width: 200,
          child: Row(
            children: [
              Flexible(
                child: Hyperlink.fromString(
                  'This is a long link that should wrap within a Flexible widget',
                  url: 'https://example.com',
                ),
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(Hyperlink), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Hyperlink unconstrained parents', () {
    testWidgets('renders in a ListView', (WidgetTester tester) async {
      await tester.pumpApp(
        ListView(
          children: [
            Hyperlink.fromString(
              'list link',
              url: 'https://example.com',
            ),
            const Text('other item'),
          ],
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('list link'), findsOneWidget);
      expect(find.text('other item'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in SingleChildScrollView > Column', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: Column(
            children: [
              Hyperlink.fromString(
                'scrollable link',
                url: 'https://example.com',
              ),
              const Text('below'),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('scrollable link'), findsOneWidget);
      expect(find.text('below'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Hyperlink inline WidgetSpan usage', () {
    testWidgets('inline renders inside Text.rich in constrained parent',
        (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(
          width: 300,
          child: Text.rich(
            TextSpan(
              text: 'Accept our ',
              children: [
                Hyperlink.inline(
                  'terms',
                  url: 'https://example.com/tos',
                ),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.byType(Hyperlink), findsOneWidget);
      expect(find.text('terms'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets(
        'inline long text in Text.rich in constrained parent is clipped',
        (WidgetTester tester) async {
      await tester.pumpApp(
        SizedBox(
          width: 50,
          height: 20,
          child: Text.rich(
            TextSpan(
              text: 'See ',
              children: [
                Hyperlink.inline(
                  'an extremely long terms of service link that will overflow',
                  url: 'https://example.com/tos',
                ),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Text.rich clips WidgetSpan children rather than causing an overflow error.
      expect(find.byType(Hyperlink), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
