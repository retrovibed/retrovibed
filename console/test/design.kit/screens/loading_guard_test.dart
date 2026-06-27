import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

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

class _Toggle extends StatefulWidget {
  const _Toggle({required this.builder});
  final Widget Function(BuildContext, VoidCallback toggle, bool visible) builder;

  @override
  _ToggleState createState() => _ToggleState();
}

class _ToggleState extends State<_Toggle> {
  bool _visible = true;

  @override
  Widget build(BuildContext context) {
    return widget.builder(context, () => setState(() => _visible = !_visible), _visible);
  }
}

void main() {
  setUp(() => _TrackedWidgetState.active.clear());

  testWidgets('no LoadingBoundary descendants never shows the overlay', (
    WidgetTester tester,
  ) async {
    await tester.pumpApp(ds.LoadingGuard(const Text('hello world')));
    await tester.pump();
    expect(find.byType(CircularProgressIndicator), findsNothing);
  });

  testWidgets('a single mounted LoadingBoundary shows the overlay', (
    WidgetTester tester,
  ) async {
    await tester.pumpApp(
      ds.LoadingGuard(ds.LoadingBoundary(const Text('hello world'))),
    );
    await tester.pump();
    await tester.pump();
    expect(find.byType(CircularProgressIndicator), findsOneWidget);
  });

  testWidgets('nested LoadingBoundary widgets show a single overlay', (
    WidgetTester tester,
  ) async {
    await tester.pumpApp(
      ds.LoadingGuard(
        ds.LoadingBoundary(ds.LoadingBoundary(const Text('hello world'))),
      ),
    );
    await tester.pump();
    await tester.pump();
    expect(find.byType(CircularProgressIndicator), findsOneWidget);
  });

  testWidgets('overlay clears only once every LoadingBoundary unmounts', (
    WidgetTester tester,
  ) async {
    late VoidCallback toggleOuter;
    late VoidCallback toggleInner;

    await tester.pumpApp(
      ds.LoadingGuard(
        _Toggle(
          builder: (context, toggle, outerVisible) {
            toggleOuter = toggle;
            if (!outerVisible) return const Text('outer done');
            return ds.LoadingBoundary(
              _Toggle(
                builder: (context, toggle, innerVisible) {
                  toggleInner = toggle;
                  if (!innerVisible) return const Text('inner done');
                  return ds.LoadingBoundary(const _TrackedWidget(tag: 'leaf'));
                },
              ),
            );
          },
        ),
      ),
    );
    await tester.pump();
    await tester.pump();
    expect(find.byType(CircularProgressIndicator), findsOneWidget);

    toggleInner();
    await tester.pump();
    await tester.pump();
    expect(find.byType(CircularProgressIndicator), findsOneWidget);

    toggleOuter();
    await tester.pump();
    await tester.pump();
    expect(find.byType(CircularProgressIndicator), findsNothing);
    expect(find.text('outer done'), findsOneWidget);
  });

  testWidgets('LoadingBoundary(loading: false) never registers as loading', (
    WidgetTester tester,
  ) async {
    await tester.pumpApp(
      ds.LoadingGuard(
        ds.LoadingBoundary(const Text('hello world'), loading: false),
      ),
    );
    await tester.pump();
    await tester.pump();
    expect(find.byType(CircularProgressIndicator), findsNothing);
  });

  testWidgets('toggling LoadingBoundary.loading updates the overlay without unmounting', (
    WidgetTester tester,
  ) async {
    late VoidCallback toggle;

    await tester.pumpApp(
      ds.LoadingGuard(
        _Toggle(
          builder: (context, t, loading) {
            toggle = t;
            return ds.LoadingBoundary(
              const _TrackedWidget(tag: 'leaf'),
              loading: loading,
            );
          },
        ),
      ),
    );
    await tester.pump();
    await tester.pump();
    expect(find.byType(CircularProgressIndicator), findsOneWidget);
    expect(_TrackedWidgetState.active, contains('leaf'));

    toggle();
    await tester.pump();
    await tester.pump();
    expect(find.byType(CircularProgressIndicator), findsNothing);
    expect(_TrackedWidgetState.active, contains('leaf'));

    toggle();
    await tester.pump();
    await tester.pump();
    expect(find.byType(CircularProgressIndicator), findsOneWidget);
    expect(_TrackedWidgetState.active, contains('leaf'));
  });
}
