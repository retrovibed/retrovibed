import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:flutter/services.dart' as services;
import './seed.dart' as seed;

class SeedTypography extends StatelessWidget {
  final String current;
  final seed.Classifier classifier;
  final void Function(seed.Seed)? onChange;

  const SeedTypography(
    this.current, {
    super.key,
    this.onChange,
    required this.classifier,
  });

  @override
  Widget build(BuildContext context) {
    final classifer = this.classifier;
    final _current = this.classifier.classify(this.current);
    if (onChange == null) {
      return _current.label;
    }

    return _SeedDropdown(
      current: _current,
      community: classifer.community,
      onChange: onChange!,
    );
  }
}

class _SeedDropdown extends StatefulWidget {
  final seed.Seed current;
  final String community;
  final void Function(seed.Seed) onChange;

  const _SeedDropdown({
    required this.current,
    required this.community,
    required this.onChange,
  });

  @override
  State<_SeedDropdown> createState() => _SeedDropdownState();
}

class _SeedDropdownState extends State<_SeedDropdown> {
  late seed.Seed _selected;
  late List<seed.Seed> _options;

  @override
  void initState() {
    super.initState();
    _selected = widget.current;
    _options = [
      seed.Seed.global(),
      // a nil community means there is no community seed to offer.
      if (!uuidx.isMin(uuidx.fromString(widget.community))) seed.Seed.community(widget.community),
      seed.Seed.unique(uuidx.random()),
      if (!uuidx.isMinMax(uuidx.fromString(widget.current.id)) && widget.current.id != widget.community)
        widget.current,
    ];
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final options = _options;

    return Container(
      width: double.infinity,
      child: Material(
        color: Colors.transparent,
        child: DropdownButton<String>(
          isDense: true,
          isExpanded: true,
          alignment: Alignment.topLeft,
          value: _selected.id,
          padding: defaults.padding.copyWith(left: 0, right: 0),
          underline: const SizedBox(),
          icon: const SizedBox(),
          onChanged: (id) {
            if (id == null) return;
            final selected = options.firstWhere(
              (s) => s.id == id,
              orElse: () => _selected,
            );
            setState(() => _selected = selected);
            widget.onChange(selected);
          },
          items: options.map((s) {
            final focused = _selected.id == s.id;
            return DropdownMenuItem<String>(
              key: ValueKey(s.id),
              value: s.id,
              child: Padding(
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  spacing: defaults.spacing,
                  children: [
                    Icon(s.icon),
                    const SizedBox(), // double the spacing to match checkboxes.
                    s.label,
                    const Spacer(),
                    if (focused) ...[
                      Expanded(
                        child: Text(
                          s.id,
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      ds.buttons.copy(
                        onPressed: () {
                          services.Clipboard.setData(
                            services.ClipboardData(text: s.id),
                          );
                        },
                        size: 12,
                      ),
                    ],
                  ],
                ),
                padding: EdgeInsetsGeometry.only(left: 9, right: 9),
              ),
            );
          }).toList(),
        ),
      ),
    );
  }
}
