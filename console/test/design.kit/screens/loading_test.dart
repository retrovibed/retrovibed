import 'package:flutter/rendering.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter/gestures.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

// Stateful widget that records whether it has been built at least once.
class _TrackedWidget extends StatefulWidget {
  const _TrackedWidget({required this.tag});
  final String tag;
  @override
  _TrackedWidgetState createState() => _TrackedWidgetState();
}

class _TrackedWidgetState extends State<_TrackedWidget> {
  static final Set<String> active = {};

  @override
  void initState() {
    super.initState();
    active.add(widget.tag);
  }

  @override
  void dispose() {
    active.remove(widget.tag);
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => const SizedBox();
}

void main() {
  group('Loading maintainState', () {
    setUp(() => _TrackedWidgetState.active.clear());

    testWidgets('maintainState:true keeps child mounted while loading', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.Loading(
          const _TrackedWidget(tag: 'child'),
          loading: true,
          maintainState: true,
        ),
      );
      await tester.pump();
      expect(_TrackedWidgetState.active, contains('child'));
    });

    testWidgets('maintainState:false unmounts child while loading', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.Loading(
          const _TrackedWidget(tag: 'child'),
          loading: true,
          maintainState: false,
        ),
      );
      await tester.pump();
      expect(_TrackedWidgetState.active, isNot(contains('child')));
    });

    testWidgets('default keeps child mounted while loading', (
      WidgetTester tester,
    ) async {
      await tester.pumpApp(
        ds.Loading(
          const _TrackedWidget(tag: 'child'),
          loading: true,
        ),
      );
      await tester.pump();
      expect(_TrackedWidgetState.active, contains('child'));
    });

    // Regression test: the overlay used to collapse to the spinner's own
    // intrinsic size while loading because Visibility only maintained the
    // child's *state*, not its layout size, so the Stack it sits in shrank
    // down as if the child had been removed from the tree.
    testWidgets('maintainState:true preserves child layout size while loading', (
      WidgetTester tester,
    ) async {
      const contentWidth = 200.0;
      const contentHeight = 120.0;
      await tester.pumpApp(
        UnconstrainedBox(
          child: ds.Loading(
            const SizedBox(width: contentWidth, height: contentHeight),
            loading: true,
            maintainState: true,
          ),
        ),
      );
      await tester.pump();

      expect(tester.getSize(find.byType(ds.Loading)), const Size(contentWidth, contentHeight));
    });

    testWidgets('maintainState:false collapses child layout size while loading', (
      WidgetTester tester,
    ) async {
      const contentWidth = 200.0;
      const contentHeight = 120.0;
      await tester.pumpApp(
        UnconstrainedBox(
          child: ds.Loading(
            const SizedBox(width: contentWidth, height: contentHeight),
            loading: true,
            maintainState: false,
          ),
        ),
      );
      await tester.pump();

      expect(tester.getSize(find.byType(ds.Loading)), isNot(const Size(contentWidth, contentHeight)));
    });
  });

  testWidgets(
    'Cursor persists over child when loading is false',
    (WidgetTester tester) async {
      final clickable = WidgetStateMouseCursor.clickable.resolve({});
      await tester.pumpApp(
        SizedBox(
          width: 200,
          height: 100,
          child: Center(
            child: ds.Loading(
              TextButton(
                child: const Text('Test Child Widget'),
                onPressed: () {},
              ),
              loading: false,
            ),
          ),
        ),
        theme: ThemeData(
          textButtonTheme: TextButtonThemeData(
            style: TextButton.styleFrom(enabledMouseCursor: clickable),
          ),
        ),
      );
      await tester.pumpAndSettle();

      // Initialize the mouse device OUTSIDE the widget area
      final gesture = await tester.createGesture(kind: PointerDeviceKind.mouse);
      await gesture.addPointer(location: Offset.zero);
      await tester.pumpAndSettle();
      expect(tester.currentCursor(), SystemMouseCursors.basic);

      final childFinder = find.text('Test Child Widget');
      expect(childFinder, findsOneWidget);

      await tester.pumpAndSettle();

      final Offset location = tester.getCenter(childFinder);
      // This triggers the 'onEnter' transition required by MouseTracker
      await gesture.moveTo(location);
      RendererBinding.instance.mouseTracker.updateAllDevices();
      await tester.pumpAndSettle();

      expect(tester.resolvedCursorAt(childFinder), clickable);
      expect(tester.currentCursor(), clickable);
    },
    variant: TargetPlatformVariant.all(),
  );
}
