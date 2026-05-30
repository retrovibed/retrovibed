import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/design.kit/theme.defaults.dart';

/// Input widget for uint64 values.
/// Treats zero and negative values (e.g. MaxUint64 stored as -1 in Int64) as
/// "unlimited" and displays an empty field. Submitting an empty field or 0
/// calls onChanged with Int64.ZERO.
class Uint64 extends StatefulWidget {
  final Int64 value;
  final ValueChanged<Int64> onChanged;
  final List<({String label, Int64 value})> presets;
  final InputDecoration decoration;

  const Uint64({
    super.key,
    required this.value,
    required this.onChanged,
    this.presets = const [],
    this.decoration = const InputDecoration(),
  });

  @override
  State<Uint64> createState() => _Uint64State();
}

class _Uint64State extends State<Uint64> {
  bool _expanded = false;
  int _generation = 0;

  String _display(Int64 v) => v > Int64.ZERO ? v.toString() : '';

  void _onTextChanged(String text) {
    final n = int.tryParse(text);
    widget.onChanged(n != null && n > 0 ? Int64(n) : Int64.ZERO);
  }

  void _selectPreset(({String label, Int64 value}) preset) {
    setState(() {
      _expanded = false;
      _generation++;
    });
    widget.onChanged(preset.value);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.start,
      spacing: defaults.spacing,
      children: [
        TextFormField(
          key: ValueKey(_generation),
          initialValue: _display(widget.value),
          keyboardType: TextInputType.number,
          inputFormatters: [FilteringTextInputFormatter.digitsOnly],
          decoration: widget.decoration.copyWith(
            isDense: true,
            hintText: widget.decoration.hintText ?? "0",
            suffixIcon:
                widget.presets.isEmpty
                    ? widget.decoration.suffixIcon
                    : IconButton(
                      onPressed: () => setState(() => _expanded = !_expanded),
                      icon: Icon(
                        _expanded ? Icons.arrow_drop_up : Icons.arrow_drop_down,
                      ),
                      padding: EdgeInsets.zero,
                      constraints: const BoxConstraints(),
                    ),
          ),
          onChanged: _onTextChanged,
        ),
        if (_expanded && widget.presets.isNotEmpty)
          LayoutBuilder(
            builder: (context, constraints) {
              final buttonWidth =
                  (constraints.maxWidth - defaults.spacing * (widget.presets.length - 1)) / widget.presets.length;
              return Wrap(
                spacing: defaults.spacing,
                runSpacing: defaults.spacing,
                children:
                    widget.presets
                        .map(
                          (p) => SizedBox(
                            width: buttonWidth,
                            child: OutlinedButton(
                              onPressed: () => _selectPreset(p),
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
                          ),
                        )
                        .toList(),
              );
            },
          ),
      ],
    );
  }
}
