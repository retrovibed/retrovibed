import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/design.kit/screens/overlay.dart' as s;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

// A widget that paints a solid color so we can find it by type and measure it.
class _ColorBox extends StatelessWidget {
  final Color color;
  const _ColorBox(this.color);

  @override
  Widget build(BuildContext context) => ColoredBox(color: color);
}

void main() {
  group('Overlay — overlay does NOT expand with StackFit.passthrough', () {
    testWidgets(
      'plain overlay widget sizes to its own content, not the stack',
      (tester) async {
        // Column gives loose constraints — same as grid.dart's Column wrapping Overlay.
        await tester.pumpApp(
          Column(
            children: [
              s.Overlay(
                const SizedBox(width: 300, height: 200, child: Text('content')),
                overlay: const SizedBox(width: 10, height: 10, child: _ColorBox(Colors.red)),
              ),
            ],
          ),
        );
        await tester.pumpAndSettle();

        final overlaySize = tester.getSize(find.byType(_ColorBox));
        // StackFit.passthrough passes loose constraints → overlay stays 10x10.
        expect(overlaySize.width, equals(10));
        expect(overlaySize.height, equals(10));
        expect(tester.takeException(), isNull);
      },
    );
  });

  group('Overlay — Positioned.fill workaround expands overlay to fill stack', () {
    testWidgets(
      'Positioned.fill overlay fills the full stack size',
      (tester) async {
        await tester.pumpApp(
          SizedBox(
            width: 300,
            height: 200,
            child: s.Overlay(
              const SizedBox(width: 300, height: 200, child: Text('content')),
              overlay: const Positioned.fill(child: _ColorBox(Colors.red)),
            ),
          ),
        );
        await tester.pumpAndSettle();

        final overlaySize = tester.getSize(find.byType(_ColorBox));
        expect(overlaySize.width, equals(300));
        expect(overlaySize.height, equals(200));
        expect(tester.takeException(), isNull);
      },
    );

    testWidgets(
      'Positioned.fill overlay fills stack when child is smaller than parent',
      (tester) async {
        // No SizedBox parent — Stack receives loose constraints and sizes to its child.
        await tester.pumpApp(
          s.Overlay(
            const SizedBox(width: 50, height: 50, child: Text('small')),
            overlay: const Positioned.fill(child: _ColorBox(Colors.blue)),
          ),
        );
        await tester.pumpAndSettle();

        final overlaySize = tester.getSize(find.byType(_ColorBox));
        // Positioned.fill fills the Stack's resolved size (50x50, set by child).
        expect(overlaySize.width, equals(50));
        expect(overlaySize.height, equals(50));
        expect(tester.takeException(), isNull);
      },
    );
  });
}
