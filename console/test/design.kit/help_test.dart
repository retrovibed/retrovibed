import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/design.kit/help.dart';
import 'package:retrovibed/designkit.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('HelpScope registration', () {
    testWidgets('Help registers description on mount', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(
            Text('child'),
            Hint(const Text('desc')),
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
                Hint(const Text('desc')),
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
              Help(Text('a'), Hint(const Text('first'))),
              Help(Text('b'), Hint(const Text('second'))),
              Help(Text('c'), Hint(const Text('third'))),
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
          Help(Text('child'), Hint(const Text('desc'))),
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
        modals.Node(
          HelpScope(
            Help(Text('child'), Hint(const Text('find things')), key: Key('help')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(of: find.byKey(Key('help')), matching: find.byType(InkWell)));
      await tester.pumpAndSettle();

      expect(find.text('find things'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('ESC closes help overlay', (tester) async {
      await tester.pumpApp(
        modals.Node(
          HelpScope(
            Help(Text('child'), Hint(const Text('desc')), key: Key('help')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(of: find.byKey(Key('help')), matching: find.byType(InkWell)));
      await tester.pumpAndSettle();
      expect(find.text('desc'), findsOneWidget);

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      await tester.pumpAndSettle();

      expect(find.text('desc'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('alt+? toggles off after second press', (tester) async {
      await tester.pumpApp(
        modals.Node(
          HelpScope(
            Help(Text('child'), Hint(const Text('desc')), key: Key('help')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(of: find.byKey(Key('help')), matching: find.byType(InkWell)));
      await tester.pumpAndSettle();
      expect(find.text('desc'), findsOneWidget);

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      expect(find.text('desc'), findsNothing);

      expect(tester.takeException(), isNull);
    });
  });

  group('HelpScope content', () {
    testWidgets('renders registered Hint widgets', (tester) async {
      await tester.pumpApp(
        modals.Node(
          HelpScope(
            Column(
              children: [
                Help(Text('a'), Hint(const Text('filter by title'))),
                Help(Text('b'), Hint(const Text('drag files')), key: Key('help-b')),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(of: find.byKey(Key('help-b')), matching: find.byType(InkWell)));
      await tester.pumpAndSettle();

      expect(find.text('drag files'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders arbitrary Widget descriptions', (tester) async {
      await tester.pumpApp(
        modals.Node(
          HelpScope(
            Help(Text('child'), Text('custom help text'), key: Key('help')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(of: find.byKey(Key('help')), matching: find.byType(InkWell)));
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
                Hint(const Text('desc')),
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
            Builder(
              builder: (context) {
                rebuildCount++;
                return Text('child');
              },
            ),
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
        modals.Node(
          HelpScope(
            Column(
              children: [
                Help(Text('a'), HelpScope.None),
                Help(Text('b'), Hint(const Text('shown')), key: Key('help-b')),
              ],
            ),
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
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(of: find.byKey(Key('help-b')), matching: find.byType(InkWell)));
      await tester.pumpAndSettle();

      expect(find.text('shown'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Help with None description does not unregister on dispose', (tester) async {
      final show = ValueNotifier(true);

      await tester.pumpApp(
        HelpScope(
          Column(
            children: [
              Help(Text('a'), Hint(const Text('desc'))),
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

  group('Help activated display', () {
    testWidgets('displays overlay when activated', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(
            Text('child'),
            Hint(const Text('desc')),
            key: Key('help'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.descendant(of: find.byKey(Key('help')), matching: find.byType(InkWell)), findsNothing);

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      scope.toggle();
      await tester.pumpAndSettle();

      expect(find.descendant(of: find.byKey(Key('help')), matching: find.byType(InkWell)), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('shows click cursor when activated', (tester) async {
      await tester.pumpApp(
        HelpScope(
          Help(
            Text('child'),
            Hint(const Text('desc')),
            key: Key('help'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      scope.toggle();
      await tester.pumpAndSettle();

      final inkWell = tester.widget<InkWell>(
        find.descendant(of: find.byKey(Key('help')), matching: find.byType(InkWell)),
      );
      expect(inkWell.mouseCursor, SystemMouseCursors.click);
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
                Hint(const Text('inner desc')),
              ),
              Hint(const Text('middle desc')),
            ),
            Hint(const Text('outer desc')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final scope = tester.state<HelpScopeState>(find.byType(HelpScope));
      expect(scope.descriptions, hasLength(3));
      expect(tester.takeException(), isNull);
    });

    testWidgets('innermost nested Help is clickable', (tester) async {
      await tester.pumpApp(
        modals.Node(
          HelpScope(
            Help(
              Help(
                Help(
                  Text('deeply wrapped'),
                  Hint(const Text('inner desc')),
                  key: Key('help-inner'),
                ),
                Hint(const Text('middle desc')),
              ),
              Hint(const Text('outer desc')),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(of: find.byKey(Key('help-inner')), matching: find.byType(InkWell)).first);
      await tester.pumpAndSettle();

      expect(find.text('inner desc'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('outer Help is clickable in area not covered by inner Help', (tester) async {
      await tester.pumpApp(
        modals.Node(
          HelpScope(
            Help(
              Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Help(
                    const SizedBox(width: 100, height: 100),
                    Hint(const Text('inner desc')),
                    key: Key('help-inner'),
                  ),
                  const SizedBox(width: 100, height: 100, key: Key('outer-only')),
                ],
              ),
              Hint(const Text('outer desc')),
              key: Key('help-outer'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();

      await tester.tapAt(tester.getCenter(find.byKey(Key('outer-only'))));
      await tester.pumpAndSettle();

      expect(find.text('outer desc'), findsOneWidget);
      expect(find.text('inner desc'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Help with tappable child shows description instead of firing child tap', (tester) async {
      var childTapped = false;
      await tester.pumpApp(
        modals.Node(
          HelpScope(
            Help(
              GestureDetector(
                onTap: () => childTapped = true,
                child: const SizedBox(width: 100, height: 100),
              ),
              Hint(const Text('desc')),
              key: Key('help'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(of: find.byKey(Key('help')), matching: find.byType(InkWell)));
      await tester.pumpAndSettle();

      expect(find.text('desc'), findsOneWidget);
      expect(childTapped, isFalse);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Help with tappable child and nested help shows description instead of firing child tap', (
      tester,
    ) async {
      var childTapped = false;
      await tester.pumpApp(
        modals.Node(
          HelpScope(
            Help(
              Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  GestureDetector(
                    onTap: () => childTapped = true,
                    child: const SizedBox(width: 100, height: 100),
                    key: Key('child'),
                  ),
                  Help(
                    const SizedBox(width: 100, height: 100, key: Key('outer-only')),
                    Hint(const Text('inner')),
                  ),
                ],
              ),
              Hint(const Text('outer')),
              key: Key('help'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.byKey(Key('child')), warnIfMissed: false);
      await tester.pumpAndSettle();

      expect(find.text('outer'), findsOneWidget);
      expect(childTapped, isFalse);
      expect(tester.takeException(), isNull);
    });

    testWidgets('inner Help with tappable child shows description instead of firing child tap', (tester) async {
      var childTapped = false;
      await tester.pumpApp(
        modals.Node(
          HelpScope(
            Help(
              Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Help(
                    GestureDetector(
                      onTap: () => childTapped = true,
                      child: const SizedBox(width: 100, height: 100),
                    ),
                    Hint(const Text('inner desc')),
                    key: Key('help-inner'),
                  ),
                  const SizedBox(width: 100, height: 100, key: Key('outer-only')),
                ],
              ),
              Hint(const Text('outer desc')),
              key: Key('help-outer'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();

      await tester.tap(find.descendant(of: find.byKey(Key('help-inner')), matching: find.byType(InkWell)).first);
      await tester.pumpAndSettle();

      expect(find.text('inner desc'), findsOneWidget);
      expect(childTapped, isFalse);
      expect(tester.takeException(), isNull);
    });
  });

  group('Help content close button', () {
    testWidgets('closing the help content does not throw after the originating widget is removed', (tester) async {
      final show = ValueNotifier(true);

      await tester.pumpApp(
        modals.Node(
          HelpScope(
            ValueListenableBuilder<bool>(
              valueListenable: show,
              builder: (_, visible, __) {
                if (!visible) return Text('removed');
                return Help(Text('child'), Hint(const Text('desc')), key: Key('help'));
              },
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.sendKeyDownEvent(LogicalKeyboardKey.altLeft);
      await tester.sendKeyEvent(LogicalKeyboardKey.slash);
      await tester.sendKeyUpEvent(LogicalKeyboardKey.altLeft);
      await tester.pumpAndSettle();
      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();
      await tester.tap(find.descendant(of: find.byKey(Key('help')), matching: find.byType(InkWell)));
      await tester.pumpAndSettle();
      expect(find.text('desc'), findsOneWidget);

      // Remove the widget that originally opened the help content while its
      // modal is still showing. The close button must not depend on that
      // widget's (now deactivated) BuildContext to find the modal node.
      show.value = false;
      await tester.pumpAndSettle();
      expect(find.text('desc'), findsOneWidget);

      await tester.tap(find.byIcon(Icons.clear));
      await tester.pumpAndSettle();

      expect(find.text('desc'), findsNothing);
      expect(tester.takeException(), isNull);
    });
  });
}
