import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/help.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('HelpScope registration', () {
    testWidgets('Help registers description on mount', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(
            Text('child'),
            Hint(label: const Text('X'), description: const Text('desc')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      expect(scope.descriptions, hasLength(1));
      expect(tester.takeException(), isNull);
    });

    testWidgets('Help unregisters description on dispose', (tester) async {
      final show = ValueNotifier(true);

      await tester.pumpApp(
        HelpScope(
          ValueListenableBuilder<bool>(
            valueListenable: show,
            builder: (_, visible, __) {
              if (!visible) return Text('empty');
              return Help(
                Text('child'),
                Hint(label: const Text('X'), description: const Text('desc')),
              );
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      expect(scope.descriptions, hasLength(1));

      show.value = false;
      await tester.pumpAndSettle();

      expect(scope.descriptions, hasLength(0));
      expect(tester.takeException(), isNull);
    });

    testWidgets('multiple Help widgets register in tree order', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Column(
            children: [
              Help(Text('a'), Hint(label: const Text('A'), description: const Text('first'))),
              Help(Text('b'), Hint(label: const Text('B'), description: const Text('second'))),
              Help(Text('c'), Hint(label: const Text('C'), description: const Text('third'))),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      expect(scope.descriptions, hasLength(3));
      expect(tester.takeException(), isNull);
    });
  });

  group('HelpScope visibility', () {
    testWidgets('overlay not visible initially', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(Text('child'), Hint(label: const Text('X'), description: const Text('desc'))),
        ),
      );
      await tester.pumpAndSettle();

      final visibility = tester.widget<Visibility>(
        find.descendant(
          of: find.byType(HelpScope),
          matching: find.byType(Visibility),
        ),
      );
      expect(visibility.visible, isFalse);
      expect(tester.takeException(), isNull);
    });

    testWidgets('alt+? toggles help overlay', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(Text('child'), Hint(label: const Text('Search'), description: const Text('find things'))),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();

      expect(find.text('Help'), findsOneWidget);
      expect(find.text('Search'), findsOneWidget);
      expect(find.text('find things'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('ESC closes help overlay', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(Text('child'), Hint(label: const Text('X'), description: const Text('desc'))),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      expect(find.text('Help'), findsOneWidget);

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      await tester.pumpAndSettle();

      expect(find.text('Help'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('alt+? toggles off after second press', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(Text('child'), Hint(label: const Text('X'), description: const Text('desc'))),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      expect(find.text('Help'), findsOneWidget);

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      expect(find.text('Help'), findsNothing);

      expect(tester.takeException(), isNull);
    });
  });

  group('HelpScope content', () {
    testWidgets('renders registered Hint widgets', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Column(
            children: [
              Help(Text('a'), Hint(label: const Text('Search'), description: const Text('filter by title'))),
              Help(Text('b'), Hint(label: const Text('Upload'), description: const Text('drag files'))),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();

      expect(find.text('Search'), findsOneWidget);
      expect(find.text('filter by title'), findsOneWidget);
      expect(find.text('Upload'), findsOneWidget);
      expect(find.text('drag files'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders arbitrary Widget descriptions', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(Text('child'), Text('custom help text')),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();

      expect(find.text('custom help text'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Help transparency', () {
    testWidgets('Help does not alter child layout', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Column(
            children: [
              Help(
                SizedBox(key: Key('target'), width: 200, height: 50),
                Hint(label: const Text('X'), description: const Text('desc')),
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      final box = tester.renderObject<RenderBox>(find.byKey(Key('target')));
      expect(box.size.width, 200.0);
      expect(box.size.height, 50.0);
      expect(tester.takeException(), isNull);
    });
  });

  group('HelpScope.of context lookup', () {
    testWidgets('returns HelpScopeState from descendant context', (tester) async {
      HelpScopeState? found;

      await tester.pumpApp(
        HelpScope(
          Builder(
            builder: (context) {
              found = HelpScope.of(context);
              return Text('child');
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(found, isNotNull);
      expect(found, isA<HelpScopeState>());
      expect(tester.takeException(), isNull);
    });
  });

  group('HelpScope.None', () {
    testWidgets('Help with None description does not subscribe to visibility changes', (tester) async {
      int rebuildCount = 0;

      await tester.pumpApp(
        HelpScope(
          Help(
            Builder(builder: (context) {
              rebuildCount++;
              return Text('child');
            }),
            HelpScope.None,
          ),
        ),
      );
      await tester.pumpAndSettle();
      rebuildCount = 0; // reset after initial build

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      scope.toggle();
      await tester.pumpAndSettle();

      expect(rebuildCount, 0);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Help with None description does not register', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(Text('child'), HelpScope.None),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      expect(scope.descriptions, hasLength(0));
      expect(tester.takeException(), isNull);
    });

    testWidgets('Help with None description does not appear in overlay', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Column(
            children: [
              Help(Text('a'), Hint(label: const Text('Visible'), description: const Text('shown'))),
              Help(Text('b'), HelpScope.None),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      expect(scope.descriptions, hasLength(1));

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();

      expect(find.text('Visible'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Help with None description does not unregister on dispose', (tester) async {
      final show = ValueNotifier(true);

      await tester.pumpApp(
        HelpScope(
          Column(
            children: [
              Help(Text('a'), Hint(label: const Text('X'), description: const Text('desc'))),
              ValueListenableBuilder<bool>(
                valueListenable: show,
                builder: (_, visible, __) {
                  if (!visible) return Text('empty');
                  return Help(Text('b'), HelpScope.None);
                },
              ),
            ],
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      expect(scope.descriptions, hasLength(1));

      show.value = false;
      await tester.pumpAndSettle();

      expect(scope.descriptions, hasLength(1));
      expect(tester.takeException(), isNull);
    });
  });

  group('Help nesting', () {
    testWidgets('nested Help wrappers register all descriptions', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(
            Help(
              Help(
                Text('deeply wrapped'),
                Hint(label: const Text('Inner'), description: const Text('inner desc')),
              ),
              Hint(label: const Text('Middle'), description: const Text('middle desc')),
            ),
            Hint(label: const Text('Outer'), description: const Text('outer desc')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      expect(scope.descriptions, hasLength(3));

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();

      expect(find.text('Inner'), findsOneWidget);
      expect(find.text('Middle'), findsOneWidget);
      expect(find.text('Outer'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });
}
