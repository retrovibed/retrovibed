import 'package:flutter/material.dart';
import './borders.dart' as borders;

class Debug extends StatelessWidget {
  final Widget? child;
  final Border? border;
  const Debug(this.child, {super.key, this.border});

  factory Debug.colored(Widget? child, {Key? key, required Color color}) {
    return Debug(child, key: key, border: Border.all(color: color));
  }

  factory Debug.green(Widget? child, {Key? key}) =>
      Debug.colored(child, key: key ?? ValueKey("green"), color: Colors.green);

  factory Debug.blue(Widget? child, {Key? key}) =>
      Debug.colored(child, key: key ?? ValueKey("blue"), color: Colors.blue);

  factory Debug.pink(Widget? child, {Key? key}) =>
      Debug.colored(child, key: key ?? ValueKey("pink"), color: Colors.pink);

  factory Debug.white(Widget? child, {Key? key}) =>
      Debug.colored(child, key: key ?? ValueKey("white"), color: Colors.white);

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        print(
          '${key} width min: ${constraints.minWidth} max: ${constraints.maxWidth}',
        );
        print(
          '${key} height min: ${constraints.minHeight} max: ${constraints.maxHeight}',
        );
        return Stack(
          fit: StackFit.passthrough,
          children: [
            child ?? const SizedBox.shrink(),
            Positioned.fill(
              child: IgnorePointer(
                child: Container(
                  decoration: BoxDecoration(border: border ?? borders.Debug),
                ),
              ),
            ),
          ],
        );
      },
    );
  }
}
