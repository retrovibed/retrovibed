import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/inputs/uint64.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

/// Simulates a real parent that feeds onChanged back as the new value,
/// causing a rebuild on every keystroke.
class _StatefulWrapper extends StatefulWidget {
  final Int64 initial;
  final List<({String label, Int64 value})> presets;

  const _StatefulWrapper({required this.initial, this.presets = const []});

  @override
  State<_StatefulWrapper> createState() => _StatefulWrapperState();
}

class _StatefulWrapperState extends State<_StatefulWrapper> {
  late Int64 _value;

  @override
  void initState() {
    super.initState();
    _value = widget.initial;
  }

  @override
  Widget build(BuildContext context) {
    return Uint64(
      value: _value,
      onChanged: (v) => setState(() => _value = v),
      presets: widget.presets,
    );
  }
}

void main() {
  group('Uint64 input', () {
    testWidgets('renders zero as empty field', (WidgetTester tester) async {
      await tester.pumpApp(Uint64(value: Int64.ZERO, onChanged: (_) {}));
      await tester.pumpAndSettle();
      expect(find.byType(TextFormField), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders positive value as text', (WidgetTester tester) async {
      await tester.pumpApp(Uint64(value: Int64(42), onChanged: (_) {}));
      await tester.pumpAndSettle();
      expect(find.text('42'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onChanged when text changes', (WidgetTester tester) async {
      Int64? changed;
      await tester.pumpApp(Uint64(value: Int64.ZERO, onChanged: (v) => changed = v));
      await tester.pumpAndSettle();
      await tester.enterText(find.byType(TextFormField), '99');
      await tester.pumpAndSettle();
      expect(changed, equals(Int64(99)));
      expect(tester.takeException(), isNull);
    });

    testWidgets('retains focus after typing causes parent rebuild', (WidgetTester tester) async {
      // Regression: key: ValueKey(widget.value) destroyed and recreated the
      // TextFormField on every keystroke, losing focus after the first character.
      await tester.pumpApp(_StatefulWrapper(initial: Int64.ZERO));
      await tester.pumpAndSettle();

      await tester.tap(find.byType(TextFormField));
      await tester.pump();

      final focusBefore = tester.binding.focusManager.primaryFocus;
      expect(focusBefore, isNotNull);

      // Enter text without re-tapping so that focus loss is detectable.
      tester.testTextInput.enterText('5');
      await tester.pump();

      expect(
        tester.binding.focusManager.primaryFocus,
        same(focusBefore),
        reason: 'text field should retain the same FocusNode after a parent rebuild triggered by onChanged',
      );
      expect(tester.takeException(), isNull);
    });

    testWidgets('no IconButton when presets list is empty', (WidgetTester tester) async {
      await tester.pumpApp(Uint64(value: Int64.ZERO, onChanged: (_) {}));
      await tester.pumpAndSettle();
      expect(find.byType(IconButton), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('preset selection updates text field display', (WidgetTester tester) async {
      final presets = [
        (label: 'small', value: Int64(100)),
        (label: 'large', value: Int64(1000)),
      ];
      await tester.pumpApp(_StatefulWrapper(initial: Int64.ZERO, presets: presets));
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();

      await tester.tap(find.text('small'));
      await tester.pumpAndSettle();

      expect(find.text('100'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('preset selection collapses the preset list', (WidgetTester tester) async {
      final presets = [(label: 'small', value: Int64(100))];
      await tester.pumpApp(_StatefulWrapper(initial: Int64.ZERO, presets: presets));
      await tester.pumpAndSettle();

      await tester.tap(find.byType(IconButton));
      await tester.pumpAndSettle();
      expect(find.text('small'), findsOneWidget);

      await tester.tap(find.text('small'));
      await tester.pumpAndSettle();
      expect(find.text('small'), findsNothing);
      expect(tester.takeException(), isNull);
    });
  });
}
