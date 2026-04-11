import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'theme.defaults.dart';

/// A widget that repeats a child widget to fill available space.
/// Uses a factory builder to generate minimum required instances.
class Repeat extends StatefulWidget {
  final Widget Function() builder;

  const Repeat(this.builder, {super.key});

  @override
  State<Repeat> createState() => _RepeatState();
}

class _RepeatState extends State<Repeat> {
  Size _detected = Size(0.0, 0.0);

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final defaults = Defaults.of(context);
        final spacing = defaults.spacing / 2;
        final viewportSize = MediaQuery.sizeOf(context);
        // Resolve the area we need to fill
        final target = constraints.constrain(viewportSize);

        int totalNeeded = 0;
        if (!_detected.isEmpty) {
          final width = (target.width / (_detected.width + spacing)).ceil();
          final height = (target.height / (_detected.height + spacing)).floor();
          totalNeeded = width * height;
        }

        return SizedBox(
          width: target.width,
          height: target.height,
          child: Stack(
            children: [
              Flow(
                delegate: _RepeatFlowDelegate(
                  childSize: _detected,
                  spacing: spacing,
                ),
                children: List.generate(totalNeeded, (_) => widget.builder()),
              ),
              Offstage(
                child: _MeasureSize(
                  onSizeMeasured: (size) {
                    if (mounted && _detected != size) {
                      setState(() => _detected = size);
                    }
                  },
                  child: ConstrainedBox(
                    constraints: BoxConstraints(
                      maxWidth: target.width - spacing,
                      maxHeight: target.height - spacing,
                    ),
                    child: widget.builder(),
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}

/// A FlowDelegate that tiles children in a grid based on a known child size.
class _RepeatFlowDelegate extends FlowDelegate {
  final Size childSize;
  final double spacing;

  _RepeatFlowDelegate({required this.childSize, required this.spacing});

  @override
  void paintChildren(FlowPaintingContext context) {
    if (childSize.isEmpty) return;

    final double width = context.size.width;
    final int horizontalCount = (width / (childSize.width + spacing)).ceil();

    // Step between items: child size + spacing (gap only between items,
    // not on the outer edges).
    final double stepX = childSize.width + spacing;
    final double stepY = childSize.height + spacing;

    for (int i = 0; i < context.childCount; i++) {
      final int xIndex = i % horizontalCount;
      final int yIndex = i ~/ horizontalCount;

      context.paintChild(
        i,
        transform: Matrix4.translationValues(
          xIndex * stepX,
          yIndex * stepY,
          0.0,
        ),
      );
    }
  }

  @override
  bool shouldRepaint(_RepeatFlowDelegate oldDelegate) {
    return childSize != oldDelegate.childSize || spacing != oldDelegate.spacing;
  }
}

/// A RenderObjectWidget used to "leak" the intrinsic size of the child up to the State.
class _MeasureSize extends SingleChildRenderObjectWidget {
  final void Function(Size size) onSizeMeasured;

  const _MeasureSize({
    required this.onSizeMeasured,
    required Widget super.child,
  });

  @override
  RenderObject createRenderObject(BuildContext context) =>
      _MeasureSizeRenderObject(onSizeMeasured);
}

class _MeasureSizeRenderObject extends RenderProxyBox {
  final void Function(Size size) onSizeMeasured;

  _MeasureSizeRenderObject(this.onSizeMeasured);

  @override
  void performLayout() {
    super.performLayout();
    if (child != null && child!.hasSize) {
      // Notify the state about the size after the layout pass completes.
      WidgetsBinding.instance.addPostFrameCallback((_) {
        onSizeMeasured(child!.size);
      });
    }
  }
}
