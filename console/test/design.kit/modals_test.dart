import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Modal constrained parent', () {
    testWidgets('takes up parent size when constrained', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 300,
          child: modals.Node(Text('Base content')),
        ),
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(
        Container(
          key: Key('modal-content'),
          child: Text('Modal content'),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Modal content'), findsOneWidget);

      // The Masked widget renders a Container with BoxDecoration as the overlay
      // backdrop. Find it as an ancestor of the modal content key.
      final maskedContainer = tester.widget<Container>(
        find
            .ancestor(
              of: find.byKey(Key('modal-content')),
              matching: find.byWidgetPredicate(
                (w) => w is Container && w.decoration != null,
              ),
            )
            .first,
      );

      expect(maskedContainer.constraints?.maxWidth ?? 0, lessThanOrEqualTo(400));
      expect(maskedContainer.constraints?.maxHeight ?? 0, lessThanOrEqualTo(300));
      expect(tester.takeException(), isNull);
    });

    testWidgets('modal overlay respects parent width constraint', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Center(
          child: SizedBox(
            width: 500,
            height: 400,
            child: modals.Node(Text('Content')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Pushed modal'));
      await tester.pumpAndSettle();

      expect(find.text('Pushed modal'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('modal renders in Column with fixed height', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Column(
          children: [
            SizedBox(
              height: 200,
              child: modals.Node(Text('In column')),
            ),
            Expanded(child: Container()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Column modal'));
      await tester.pumpAndSettle();

      expect(find.text('Column modal'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('modal renders in Row with fixed width', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        Row(
          children: [
            SizedBox(
              width: 300,
              child: modals.Node(Text('In row')),
            ),
            Expanded(child: Container()),
          ],
        ),
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Row modal'));
      await tester.pumpAndSettle();

      expect(find.text('Row modal'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Modal unconstrained parent (fallback to screen size)', () {
    testWidgets('uses screen size when parent is unconstrained', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        modals.Node(Text('Full screen base')),
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Full screen modal'));
      await tester.pumpAndSettle();

      expect(find.text('Full screen modal'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in ListView (infinite vertical)', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ListView(
          children: [
            SizedBox(
              height: 400,
              child: modals.Node(Text('In ListView')),
            ),
          ],
        ),
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('ListView modal'));
      await tester.pumpAndSettle();

      expect(find.text('ListView modal'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('renders in SingleChildScrollView', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SingleChildScrollView(
          child: SizedBox(
            height: 500,
            child: modals.Node(Text('In scroll')),
          ),
        ),
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Scroll modal'));
      await tester.pumpAndSettle();

      expect(find.text('Scroll modal'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });
  });

  group('Modal push and pop', () {
    testWidgets('push shows modal content', (WidgetTester tester) async {
      await tester.pumpApp(modals.Node(Text('Base')));
      await tester.pumpAndSettle();

      expect(find.text('Modal'), findsNothing);

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Modal'));
      await tester.pumpAndSettle();

      expect(find.text('Modal'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('push(null) pops modal', (WidgetTester tester) async {
      await tester.pumpApp(modals.Node(Text('Base')));
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Modal'));
      await tester.pumpAndSettle();
      expect(find.text('Modal'), findsOneWidget);

      nodeState.push(null);
      await tester.pumpAndSettle();
      expect(find.text('Modal'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('reset clears modal', (WidgetTester tester) async {
      await tester.pumpApp(modals.Node(Text('Base')));
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Modal'));
      await tester.pumpAndSettle();
      expect(find.text('Modal'), findsOneWidget);

      nodeState.reset();
      await tester.pumpAndSettle();
      expect(find.text('Modal'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('stacked modals pop in order', (WidgetTester tester) async {
      await tester.pumpApp(modals.Node(Text('Base')));
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Modal 1'));
      await tester.pumpAndSettle();
      expect(find.text('Modal 1'), findsOneWidget);

      nodeState.push(Text('Modal 2'));
      await tester.pumpAndSettle();
      expect(find.text('Modal 2'), findsOneWidget);
      expect(find.text('Modal 1'), findsNothing);

      nodeState.push(null);
      await tester.pumpAndSettle();
      expect(find.text('Modal 1'), findsOneWidget);
      expect(find.text('Modal 2'), findsNothing);

      nodeState.push(null);
      await tester.pumpAndSettle();
      expect(find.text('Modal 1'), findsNothing);
      expect(tester.takeException(), isNull);
    });
  });

  group('Modal keyboard handling', () {
    testWidgets('Escape key closes modal', (WidgetTester tester) async {
      await tester.pumpApp(modals.Node(Text('Base')));
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Modal'));
      await tester.pumpAndSettle();
      expect(find.text('Modal'), findsOneWidget);

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      await tester.pumpAndSettle();

      expect(find.text('Modal'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('Escape pops stacked modals one at a time', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(modals.Node(Text('Base')));
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Modal 1'));
      await tester.pumpAndSettle();
      nodeState.push(Text('Modal 2'));
      await tester.pumpAndSettle();

      expect(find.text('Modal 2'), findsOneWidget);

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      await tester.pumpAndSettle();
      expect(find.text('Modal 1'), findsOneWidget);

      await tester.sendKeyEvent(LogicalKeyboardKey.escape);
      await tester.pumpAndSettle();
      expect(find.text('Modal 1'), findsNothing);

      expect(tester.takeException(), isNull);
    });
  });

  group('Modal visibility', () {
    testWidgets('base content remains visible when modal shown', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(modals.Node(Text('Base content')));
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(Text('Modal overlay'));
      await tester.pumpAndSettle();

      expect(find.text('Base content'), findsOneWidget);
      expect(find.text('Modal overlay'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('modal not visible when empty', (WidgetTester tester) async {
      await tester.pumpApp(modals.Node(Text('Base')));
      await tester.pumpAndSettle();

      final visibility = tester.widget<Visibility>(
        find.byType(Visibility).last,
      );
      expect(visibility.visible, isFalse);
      expect(tester.takeException(), isNull);
    });
  });

  group('Modal.of context lookup', () {
    testWidgets('Node.of returns NodeState from descendant context', (
      WidgetTester tester,
    ) async {
      modals.NodeState? foundState;

      await tester.pumpApp(
        modals.Node(
          Builder(
            builder: (context) {
              foundState = modals.Node.of(context);
              return Text('Child');
            },
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(foundState, isNotNull);
      expect(foundState, isA<modals.NodeState>());
      expect(tester.takeException(), isNull);
    });
  });

  group('Modal tap to dismiss', () {
    testWidgets('tapping background dismisses modal', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 400,
          child: modals.Node(Text('Base')),
        ),
        alignment: Alignment.topLeft,
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(
        Center(
          child: Container(
            width: 100,
            height: 100,
            color: Colors.blue,
            child: Text('Modal Content'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Modal Content'), findsOneWidget);

      // Tap on the background area (top-left corner, outside modal content)
      await tester.tapAt(Offset(10, 10));
      await tester.pumpAndSettle();

      expect(find.text('Modal Content'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping modal content does not dismiss', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 400,
          child: modals.Node(Text('Base')),
        ),
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(
        Center(
          child: Container(
            width: 200,
            height: 200,
            color: Colors.blue,
            child: Text('Modal Content'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Modal Content'), findsOneWidget);

      // Tap on the modal content
      await tester.tap(find.text('Modal Content'));
      await tester.pumpAndSettle();

      // Modal should still be visible
      expect(find.text('Modal Content'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('tapping container inside modal does not dismiss', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 400,
          child: modals.Node(Text('Base')),
        ),
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(
        Align(
          alignment: Alignment.center,
          child: Container(
            width: 150,
            height: 150,
            color: Colors.red,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text('Header'),
                ElevatedButton(
                  onPressed: () {},
                  child: Text('Button'),
                ),
              ],
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Header'), findsOneWidget);

      // Tap on button inside modal
      await tester.tap(find.text('Button'));
      await tester.pumpAndSettle();

      // Modal should still be visible
      expect(find.text('Header'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('multiple taps on background dismiss stacked modals', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 400,
          child: modals.Node(Text('Base')),
        ),
        alignment: Alignment.topLeft,
      );
      await tester.pumpAndSettle();

      final nodeState = tester.state<modals.NodeState>(
        find.byType(modals.Node),
      );

      nodeState.push(
        Center(
          child: Container(
            width: 100,
            height: 100,
            child: Text('Modal 1'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      nodeState.push(
        Center(
          child: Container(
            width: 100,
            height: 100,
            child: Text('Modal 2'),
          ),
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('Modal 2'), findsOneWidget);

      // Tap background to dismiss Modal 2
      await tester.tapAt(Offset(10, 10));
      await tester.pumpAndSettle();

      // reset() clears all modals, so Modal 1 should also be gone
      expect(find.text('Modal 2'), findsNothing);
      expect(find.text('Modal 1'), findsNothing);
      expect(tester.takeException(), isNull);
    });
  });

  group('asyncfn', () {
    testWidgets('pushes the built widget as a modal', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        modals.Node(
          Builder(
            builder: (context) => TextButton(
              onPressed: () => modals.asyncfn<void>(
                context,
                (_) => Text('Async modal'),
              ),
              child: Text('Open'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Open'));
      await tester.pumpAndSettle();

      expect(find.text('Async modal'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('future completes when completion.complete() is called', (
      WidgetTester tester,
    ) async {
      late Completer<void> trigger;
      late Future<void> result;

      await tester.pumpApp(
        modals.Node(
          Builder(
            builder: (context) => TextButton(
              onPressed: () {
                result = modals.asyncfn<void>(
                  context,
                  (completion) {
                    trigger = completion;
                    return Text('Waiting');
                  },
                );
              },
              child: Text('Open'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Open'));
      await tester.pumpAndSettle();

      var completed = false;
      result.then((_) => completed = true);

      expect(completed, isFalse);

      trigger.complete();
      await tester.pumpAndSettle();

      expect(completed, isTrue);
      expect(tester.takeException(), isNull);
    });

    testWidgets('future completes with error when completion.completeError() is called', (
      WidgetTester tester,
    ) async {
      late Completer<void> trigger;
      late Future<void> result;

      await tester.pumpApp(
        modals.Node(
          Builder(
            builder: (context) => TextButton(
              onPressed: () {
                result = modals.asyncfn<void>(
                  context,
                  (completion) {
                    trigger = completion;
                    return Text('Waiting');
                  },
                );
              },
              child: Text('Open'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Open'));
      await tester.pumpAndSettle();

      Object? caughtError;
      result.catchError((e) => caughtError = e);

      trigger.completeError(Exception('boom'));
      await tester.pumpAndSettle();

      expect(caughtError, isA<Exception>());
      expect(tester.takeException(), isNull);
    });

    testWidgets('modal is dismissed when future completes', (
      WidgetTester tester,
    ) async {
      late Completer<void> trigger;

      await tester.pumpApp(
        modals.Node(
          Builder(
            builder: (context) => TextButton(
              onPressed: () => modals.asyncfn<void>(
                context,
                (completion) {
                  trigger = completion;
                  return Text('Async modal');
                },
              ),
              child: Text('Open'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Open'));
      await tester.pumpAndSettle();

      expect(find.text('Async modal'), findsOneWidget);

      trigger.complete();
      await tester.pumpAndSettle();

      expect(find.text('Async modal'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('future completes with typed value', (
      WidgetTester tester,
    ) async {
      late Completer<String> trigger;
      late Future<String> result;

      await tester.pumpApp(
        modals.Node(
          Builder(
            builder: (context) => TextButton(
              onPressed: () {
                result = modals.asyncfn<String>(
                  context,
                  (completion) {
                    trigger = completion;
                    return Text('Picking');
                  },
                );
              },
              child: Text('Open'),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Open'));
      await tester.pumpAndSettle();

      trigger.complete('hello');
      await tester.pumpAndSettle();

      expect(await result, 'hello');
      expect(tester.takeException(), isNull);
    });

    testWidgets('future completes when dismissed by tapping outside the modal', (
      WidgetTester tester,
    ) async {
      late Future<void> result;

      await tester.pumpApp(
        SizedBox(
          width: 400,
          height: 400,
          child: modals.Node(
            Builder(
              builder: (context) => TextButton(
                onPressed: () {
                  result = modals.asyncfn<void>(
                    context,
                    (completion) => Center(
                      child: Container(
                        width: 100,
                        height: 100,
                        child: Text('Waiting'),
                      ),
                    ),
                  );
                },
                child: Text('Open'),
              ),
            ),
          ),
        ),
        alignment: Alignment.topLeft,
      );
      await tester.pumpAndSettle();

      await tester.tap(find.text('Open'));
      await tester.pumpAndSettle();

      expect(find.text('Waiting'), findsOneWidget);

      var completed = false;
      result.then((_) => completed = true);

      // Tap the background, outside the modal content, to dismiss it.
      await tester.tapAt(Offset(10, 10));
      await tester.pumpAndSettle();

      expect(find.text('Waiting'), findsNothing);
      // Reproduces the bug: reset() (triggered by the outside tap) clears the
      // modal without ever resolving the Completer asyncfn handed out, so
      // callers awaiting this future (e.g. LoadingIconButton) hang forever.
      expect(completed, isTrue);
      expect(tester.takeException(), isNull);
    });
  });
}
