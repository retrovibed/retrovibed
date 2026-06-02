import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/inputs/mimetype.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

class _StatefulWrapper extends StatefulWidget {
  final String initial;
  final List<MimetypePreset> presets;

  const _StatefulWrapper({required this.initial, this.presets = const []});

  @override
  State<_StatefulWrapper> createState() => _StatefulWrapperState();
}

class _StatefulWrapperState extends State<_StatefulWrapper> {
  late String _value;

  @override
  void initState() {
    super.initState();
    _value = widget.initial;
  }

  @override
  Widget build(BuildContext context) {
    return Mimetype(
      value: _value,
      onChanged: (v) => setState(() => _value = v),
      presets: widget.presets,
    );
  }
}

void main() {
  group('Mimetype input', () {
    final presets = [
      (label: Text('video') as Widget, value: 'video/*'),
      (label: Text('audio') as Widget, value: 'audio/*'),
      (label: Text('image') as Widget, value: 'image/*'),
    ];

    group('screen resolutions', () {
      testWidgets('renders at minimum width (300x568)', (WidgetTester tester) async {
        tester.view.physicalSize = Size(300, 568);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpApp(
          Mimetype(value: '', onChanged: (_) {}, presets: presets),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders on small mobile (320x568)', (WidgetTester tester) async {
        tester.view.physicalSize = Size(320, 568);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpApp(
          Mimetype(value: '', onChanged: (_) {}, presets: presets),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(find.byType(IconButton), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders on iPhone SE (375x667)', (WidgetTester tester) async {
        tester.view.physicalSize = Size(375, 667);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpApp(
          Mimetype(value: '', onChanged: (_) {}, presets: presets),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(find.byType(IconButton), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders on tablet (768x1024)', (WidgetTester tester) async {
        tester.view.physicalSize = Size(768, 1024);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpApp(
          Mimetype(value: '', onChanged: (_) {}, presets: presets),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(find.byType(IconButton), findsOneWidget);
        expect(tester.takeException(), isNull);
      });

      testWidgets('renders on desktop (1920x1080)', (WidgetTester tester) async {
        tester.view.physicalSize = Size(1920, 1080);
        tester.view.devicePixelRatio = 1.0;
        addTearDown(() => tester.view.reset());

        await tester.pumpApp(
          Mimetype(value: '', onChanged: (_) {}, presets: presets),
        );
        await tester.pumpAndSettle();

        expect(find.byType(TextFormField), findsOneWidget);
        expect(find.byType(IconButton), findsOneWidget);
        expect(tester.takeException(), isNull);
      });
    });

    testWidgets('renders current value in text field', (WidgetTester tester) async {
      await tester.pumpApp(
        Mimetype(value: 'video/mp4', onChanged: (_) {}),
      );
      await tester.pumpAndSettle();

      expect(find.text('video/mp4'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders empty field when value is empty', (WidgetTester tester) async {
      await tester.pumpApp(
        Mimetype(value: '', onChanged: (_) {}),
      );
      await tester.pumpAndSettle();

      expect(find.byType(TextFormField), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('no IconButton when presets list is empty', (WidgetTester tester) async {
      await tester.pumpApp(
        Mimetype(value: '', onChanged: (_) {}),
      );
      await tester.pumpAndSettle();

      expect(find.byType(IconButton), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onChanged when text changes', (WidgetTester tester) async {
      String? changed;

      await tester.pumpApp(
        Mimetype(value: '', onChanged: (v) => changed = v),
      );
      await tester.pumpAndSettle();

      await tester.enterText(find.byType(TextFormField), 'audio/flac');
      await tester.pumpAndSettle();

      expect(changed, equals('audio/flac'));
      expect(tester.takeException(), isNull);
    });

    testWidgets('dropdown button toggles expanded state', (WidgetTester tester) async {
      await tester.pumpApp(
        Mimetype(value: '', onChanged: (_) {}, presets: presets),
      );
      await tester.pumpAndSettle();

      expect(find.text('presets'), findsNothing);

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      expect(find.text('presets'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('displays preset buttons when expanded', (WidgetTester tester) async {
      await tester.pumpApp(
        Mimetype(value: '', onChanged: (_) {}, presets: presets),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      expect(find.text('video'), findsOneWidget);
      expect(find.text('audio'), findsOneWidget);
      expect(find.text('image'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('preset selection calls onChanged with preset value', (WidgetTester tester) async {
      String? changed;

      await tester.pumpApp(
        Mimetype(value: '', onChanged: (v) => changed = v, presets: presets),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      await tester.tap(find.text('audio'));
      await tester.pumpAndSettle();

      expect(changed, equals('audio/*'));
      expect(tester.takeException(), isNull);
    });

    testWidgets('preset selection updates text field display', (WidgetTester tester) async {
      await tester.pumpApp(
        _StatefulWrapper(initial: '', presets: presets),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      await tester.tap(find.text('video'));
      await tester.pumpAndSettle();

      expect(find.text('video/*'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('preset selection collapses the preset list', (WidgetTester tester) async {
      await tester.pumpApp(
        _StatefulWrapper(initial: '', presets: presets),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();
      expect(find.text('video'), findsOneWidget);

      await tester.tap(find.text('video'));
      await tester.pumpAndSettle();
      expect(find.text('video'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders custom Widget label in preset button', (WidgetTester tester) async {
      final iconPresets = [
        (label: Icon(Icons.movie) as Widget, value: 'video/*'),
        (label: Icon(Icons.music_note) as Widget, value: 'audio/*'),
      ];

      await tester.pumpApp(
        Mimetype(value: '', onChanged: (_) {}, presets: iconPresets),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.movie), findsOneWidget);
      expect(find.byIcon(Icons.music_note), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('retains focus after typing causes parent rebuild', (WidgetTester tester) async {
      await tester.pumpApp(_StatefulWrapper(initial: ''));
      await tester.pumpAndSettle();

      await tester.tap(find.byType(TextFormField));
      await tester.pump();

      final focusBefore = tester.binding.focusManager.primaryFocus;
      expect(focusBefore, isNotNull);

      tester.testTextInput.enterText('v');
      await tester.pump();

      expect(
        tester.binding.focusManager.primaryFocus,
        same(focusBefore),
        reason: 'text field should retain the same FocusNode after a parent rebuild triggered by onChanged',
      );
      expect(tester.takeException(), isNull);
    });
  });
}
