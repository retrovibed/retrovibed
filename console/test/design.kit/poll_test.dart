import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/poll.dart';
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('Poll', () {
    testWidgets('renders child immediately', (tester) async {
      await tester.pumpApp(
        Poll(
          const Text('child'),
          interval: Duration.zero,
          onTick: () async => false,
        ),
      );
      await tester.pumpAndSettle();

      expect(find.text('child'), findsOneWidget);
      expect(tester.takeException(), isNull);
    });

    testWidgets('does not tick when interval is Duration.zero', (tester) async {
      var ticks = 0;

      await tester.pumpApp(
        Poll(
          const Text('child'),
          interval: Duration.zero,
          onTick: () async {
            ticks++;
            return false;
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.pump(const Duration(seconds: 10));

      expect(ticks, 0);
      expect(tester.takeException(), isNull);
    });

    testWidgets('calls onTick on every interval while not done', (tester) async {
      var ticks = 0;

      await tester.pumpApp(
        Poll(
          const Text('child'),
          interval: const Duration(seconds: 1),
          onTick: () async {
            ticks++;
            return false;
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.pump(const Duration(seconds: 1));
      await tester.pump();
      expect(ticks, 1);

      await tester.pump(const Duration(seconds: 1));
      await tester.pump();
      expect(ticks, 2);

      await tester.pump(const Duration(seconds: 1));
      await tester.pump();
      expect(ticks, 3);
      expect(tester.takeException(), isNull);
    });

    testWidgets('stops ticking once onTick resolves true', (tester) async {
      var ticks = 0;

      await tester.pumpApp(
        Poll(
          const Text('child'),
          interval: const Duration(seconds: 1),
          onTick: () async {
            ticks++;
            return ticks >= 2;
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.pump(const Duration(seconds: 1));
      await tester.pump();
      expect(ticks, 1);

      await tester.pump(const Duration(seconds: 1));
      await tester.pump();
      expect(ticks, 2);

      await tester.pump(const Duration(seconds: 1));
      await tester.pump();
      expect(ticks, 2);

      await tester.pump(const Duration(seconds: 1));
      await tester.pump();
      expect(ticks, 2);
      expect(tester.takeException(), isNull);
    });

    testWidgets('switching interval from zero to non-zero enables polling', (tester) async {
      var ticks = 0;
      var interval = Duration.zero;
      late StateSetter setLocalState;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            setLocalState = setState;
            return Poll(
              const Text('child'),
              interval: interval,
              onTick: () async {
                ticks++;
                return false;
              },
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.pump(const Duration(seconds: 5));
      expect(ticks, 0);

      setLocalState(() => interval = const Duration(seconds: 1));
      await tester.pump();

      await tester.pump(const Duration(seconds: 1));
      await tester.pump();
      expect(ticks, 1);
      expect(tester.takeException(), isNull);
    });

    testWidgets('switching interval to zero cancels the timer', (tester) async {
      var ticks = 0;
      var interval = const Duration(seconds: 1);
      late StateSetter setLocalState;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            setLocalState = setState;
            return Poll(
              const Text('child'),
              interval: interval,
              onTick: () async {
                ticks++;
                return false;
              },
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.pump(const Duration(seconds: 1));
      await tester.pump();
      expect(ticks, 1);

      setLocalState(() => interval = Duration.zero);
      await tester.pump();

      await tester.pump(const Duration(seconds: 10));
      expect(ticks, 1);
      expect(tester.takeException(), isNull);
    });

    testWidgets('cancels the timer on dispose', (tester) async {
      var ticks = 0;
      var mounted = true;
      late StateSetter setLocalState;

      await tester.pumpApp(
        StatefulBuilder(
          builder: (context, setState) {
            setLocalState = setState;
            if (!mounted) return const SizedBox.shrink();
            return Poll(
              const Text('child'),
              interval: const Duration(seconds: 1),
              onTick: () async {
                ticks++;
                return false;
              },
            );
          },
        ),
      );
      await tester.pumpAndSettle();

      await tester.pump(const Duration(seconds: 1));
      await tester.pump();
      expect(ticks, 1);

      setLocalState(() => mounted = false);
      await tester.pump();

      await tester.pump(const Duration(seconds: 10));
      expect(ticks, 1);
      expect(tester.takeException(), isNull);
    });
  });
}
