import 'package:flutter/material.dart';
import '../theme.defaults.dart';

class Presets<T> extends StatelessWidget {
  final T current;
  final List<({String label, T Function(T current) apply})> presets;
  final ValueChanged<T> onSelected;

  const Presets({
    super.key,
    required this.current,
    required this.presets,
    required this.onSelected,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    return Wrap(
      spacing: defaults.spacing,
      runSpacing: defaults.spacing,
      alignment: WrapAlignment.start,
      children: presets
          .map(
            (p) => OutlinedButton(
              onPressed: () => onSelected(p.apply(current)),
              style: OutlinedButton.styleFrom(
                padding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 8,
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
              child: Text(p.label),
            ),
          )
          .toList(),
    );
  }
}
