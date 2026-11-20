import 'package:material_ui/material_ui.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/discovery.dart' as discovery;
import 'package:retrovibed/library/network.settings.dart';
import 'package:retrovibed/media.dart' as media;

class _GridItem extends StatelessWidget {
  final Widget child;
  final List<Widget> leading;
  final List<Widget> trailing;
  const _GridItem(
    this.child, {
    this.leading = const [],
    this.trailing = const [],
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Card(
      alignment: Alignment.topLeft,
      leading: leading,
      ClipRRect(borderRadius: defaults.borderRadius, child: child),
      trailing: trailing,
    );
  }
}

class GridSettings extends StatelessWidget {
  final media.DownloadSearchRequest current;
  final void Function(media.DownloadSearchRequest r) onChange;
  const GridSettings(this.current, {super.key, required this.onChange});

  @override
  Widget build(BuildContext context) {
    final items = [
      _GridItem(NetworkSettings()),
      _GridItem(discovery.Settings()),
    ];

    return LayoutBuilder(
      builder: (context, constraints) {
        const minheight = 812.0;
        const minCrossAxisExtent = 400.0;
        const maxCrossAxisExtent = 508.0;

        // returns (width, height) for a grid item at the given max extent,
        // where height is always minheight.
        final calc = (double extent) {
          final columns = (constraints.maxWidth / extent).ceil();
          final w = constraints.maxWidth / columns;
          return (w, minheight);
        };

        final (minW, minH) = calc(minCrossAxisExtent);
        final (maxW, maxH) = calc(maxCrossAxisExtent);
        final maxRatio = maxW / maxH;
        final compact = (minW > minH) && (maxW > maxH);

        if (compact) {
          return Column(mainAxisSize: MainAxisSize.min, children: items);
        }

        return ds.Grid(
          (context, i) => i,
          padding: EdgeInsets.zero,
          maxCrossAxisExtent: maxCrossAxisExtent,
          aspectRatio: maxRatio,
          children: items,
        );
      },
    );
  }
}
