import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

class _Probe extends StatefulWidget {
  final void Function() onInit;

  const _Probe({super.key, required this.onInit});

  @override
  State<_Probe> createState() => _ProbeState();
}

class _ProbeState extends State<_Probe> {
  int _counter = 0;

  void bump() => setState(() => _counter++);

  @override
  void initState() {
    super.initState();
    widget.onInit();
  }

  @override
  Widget build(BuildContext context) => Text('probe-$_counter');
}

void main() {
  group('LazyVisible', () {
    testWidgets('child State is never created before visible first becomes true', (tester) async {
      var inits = 0;

      await tester.pumpApp(
        ds.LazyVisible(
          _Probe(onInit: () => inits++),
          visible: false,
        ),
      );
      await tester.pumpAndSettle();

      expect(inits, 0);
      expect(find.byType(_Probe), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('child mounts and initState runs exactly once when visible turns true', (tester) async {
      var inits = 0;
      bool visible = false;
      late StateSetter setLocalState;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            setLocalState = setState;
            return ds.LazyVisible(
              _Probe(onInit: () => inits++),
              visible: visible,
            );
          },
        ),
      );
      await tester.pumpAndSettle();
      expect(inits, 0);

      setLocalState(() => visible = true);
      await tester.pumpAndSettle();

      expect(inits, 1);
      expect(find.text('probe-0'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('State persists across a hide/show cycle when maintainState is true', (tester) async {
      bool visible = true;
      late StateSetter setLocalState;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            setLocalState = setState;
            return ds.LazyVisible(
              _Probe(key: const ValueKey('probe'), onInit: () {}),
              visible: visible,
              maintainState: true,
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      final before = tester.state<_ProbeState>(find.byType(_Probe));
      before.bump();
      await tester.pump();
      expect(find.text('probe-1'), findsOneWidget);

      setLocalState(() => visible = false);
      await tester.pumpAndSettle();
      expect(find.text('probe-1'), findsNothing);

      setLocalState(() => visible = true);
      await tester.pumpAndSettle();

      final after = tester.state<_ProbeState>(find.byType(_Probe));
      expect(identical(before, after), isTrue);
      expect(find.text('probe-1'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('a differently-configured child passed while still visible is picked up immediately', (
      tester,
    ) async {
      String label = 'first';
      late StateSetter setLocalState;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            setLocalState = setState;
            return ds.LazyVisible(
              Text(label),
              visible: true,
            );
          },
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('first'), findsOneWidget);

      setLocalState(() => label = 'second');
      await tester.pumpAndSettle();

      expect(find.text('second'), findsOneWidget);
      expect(find.text('first'), findsNothing);
      expect(tester.takeException(), isNull);
    });

    testWidgets('unrelated ancestor rebuilds while never visible still do not mount the child', (tester) async {
      var inits = 0;
      bool tick = false;
      late StateSetter setLocalState;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            setLocalState = setState;
            return Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text('tick-$tick'),
                ds.LazyVisible(
                  _Probe(onInit: () => inits++),
                  visible: false,
                ),
              ],
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      setLocalState(() => tick = true);
      await tester.pumpAndSettle();
      setLocalState(() => tick = false);
      await tester.pumpAndSettle();

      expect(inits, 0);
      expect(tester.takeException(), isNull);
    });
  });
}
