import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'theme.defaults.dart';
import 'screens/loading.dart';
import 'container.dart' as _c;
import 'errors.dart';

class CarouselRow extends StatefulWidget {
  final BoxConstraints? constraints;
  final Widget title;
  final List<Widget> items;
  final bool loading;
  final Widget cause;
  final Widget background;
  final Widget empty;

  const CarouselRow({
    super.key,
    required this.title,
    required this.items,
    this.loading = false,
    this.cause = Error.zero,
    this.background = const SizedBox(),
    this.empty = const SizedBox(),
    this.constraints,
  });

  @override
  State<CarouselRow> createState() => _CarouselRowState();
}

class _CarouselRowState extends State<CarouselRow> {
  final _scroll = ScrollController();
  final _focus = FocusNode();

  @override
  void dispose() {
    _scroll.dispose();
    _focus.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    const _scrollStep = 200.0;
    final defaults = Defaults.of(context);
    return Focus(
      focusNode: _focus,
      onKeyEvent: (node, event) {
        if (!(event is KeyDownEvent || event is KeyRepeatEvent)) {
          return KeyEventResult.ignored;
        }

        if (event.logicalKey != LogicalKeyboardKey.arrowRight && event.logicalKey != LogicalKeyboardKey.arrowLeft) {
          return KeyEventResult.ignored;
        }

        double delta = _scrollStep;
        if (event.logicalKey == LogicalKeyboardKey.arrowLeft) {
          delta = delta * -1.0;
        }

        _scroll.jumpTo(
          (_scroll.offset + delta).clamp(
            _scroll.position.minScrollExtent,
            _scroll.position.maxScrollExtent,
          ),
        );

        return KeyEventResult.handled;
      },
      onFocusChange: (_) => setState(() {}),
      child: GestureDetector(
        onTap: _focus.requestFocus,
        child: _c.Container(
          constraints: widget.constraints,
          padding: defaults.padding,
          Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            spacing: defaults.spacing,
            children: [
              widget.title,
              Expanded(
                child: _c.Container(
                  clipBehavior: Clip.antiAlias,
                  decoration: BoxDecoration(
                    color: Colors.transparent,
                    borderRadius: defaults.borderRadius,
                    boxShadow: _focus.hasFocus ? defaults.highlightTint : [],
                  ),
                  Loading(
                    loading: widget.loading,
                    cause: widget.cause,
                    Stack(
                      children: [
                        if (widget.items.isEmpty) widget.background,
                        if (widget.items.isEmpty && !widget.loading) Center(child: widget.empty),
                        Listener(
                          onPointerSignal: (event) {
                            if (event is PointerScrollEvent) {
                              final delta = event.scrollDelta.dx + event.scrollDelta.dy;
                              _scroll.jumpTo(
                                (_scroll.offset + delta).clamp(
                                  _scroll.position.minScrollExtent,
                                  _scroll.position.maxScrollExtent,
                                ),
                              );
                            }
                          },
                          onPointerMove: (event) {
                            _scroll.jumpTo(
                              (_scroll.offset - event.delta.dx).clamp(
                                _scroll.position.minScrollExtent,
                                _scroll.position.maxScrollExtent,
                              ),
                            );
                          },
                          child: SingleChildScrollView(
                            controller: _scroll,
                            scrollDirection: Axis.horizontal,
                            physics: const NeverScrollableScrollPhysics(),
                            child: Row(
                              mainAxisSize: MainAxisSize.min,
                              spacing: defaults.spacing,
                              children: widget.items,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
