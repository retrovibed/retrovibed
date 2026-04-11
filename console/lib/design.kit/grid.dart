import 'package:flutter/material.dart';
import 'package:fixnum/fixnum.dart' as fixnum;
import './errors.dart' as errors;
import './help.dart';
import './theme.defaults.dart';
import './screens.dart' as screens;

class Grid<T> extends StatelessWidget {
  static fixnum.Int64 int64(int n) => fixnum.Int64(n);

  final Widget Function(BuildContext context, T i) render;
  final List<T> children;
  final Widget empty;
  final Widget overlay;
  final List<Widget> leading;
  final List<Widget> trailing;
  final bool loading;
  final Widget cause;
  final double maxCrossAxisExtent;
  final double aspectRatio;
  final EdgeInsets? padding;
  final ScrollPhysics? physics;
  final Widget help;

  const Grid(
    this.render, {
    super.key,
    this.leading = const [],
    this.trailing = const [],
    this.children = const [],
    this.loading = false,
    this.cause = errors.Error.zero,
    this.empty = const SizedBox(),
    this.overlay = const SizedBox(),
    this.maxCrossAxisExtent = 512,
    this.aspectRatio = 2 / 3,
    this.padding,
    this.physics,
    this.help = HelpScope.None,
  });

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final defaults = Defaults.of(context);
        final bounded = constraints.maxHeight.isFinite;
        final compact = defaults.isCompact;
        final grid = GridView.builder(
          shrinkWrap: true,
          reverse: compact,
          physics: physics ?? NeverScrollableScrollPhysics(),
          padding: padding ?? defaults.padding,
          itemCount: this.children.length,
          gridDelegate: SliverGridDelegateWithMaxCrossAxisExtent(
            maxCrossAxisExtent: maxCrossAxisExtent,
            crossAxisSpacing: defaults.spacing / 2, // Spacing between columns
            mainAxisSpacing: defaults.spacing / 2, // Spacing between rows
            childAspectRatio: aspectRatio,
          ),
          itemBuilder: (ctx, index) {
            return this.render(ctx, this.children.elementAt(index));
          },
        );

        final content = screens.Loading(
          screens.Overlay(
            children.length == 0 ? empty : grid,
            overlay: overlay,
          ),
          loading: loading,
          cause: cause,
        );

        return Help(
          Column(
            mainAxisAlignment: MainAxisAlignment.start,
            mainAxisSize: MainAxisSize.max,
            verticalDirection: compact ? VerticalDirection.up : VerticalDirection.down,
            spacing: 0,
            children: [
              ...leading,
              bounded ? Expanded(child: content) : content,
              ...trailing,
            ],
          ),
          help,
        );
      },
    );
  }
}
