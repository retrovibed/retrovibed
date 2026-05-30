import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import './../../design.kit/theme.defaults.dart';

/// RateLimit input for configuring rate limits.
///
/// NOTE: Units are not yet fully implemented. The units parameter should be
/// updated to use Duration objects instead of String labels for better type
/// safety and clearer semantics.
typedef RateLimitPreset = ({String label, int value, String unit});

const _defaultUnits = ['sec'];

class RateLimit extends StatefulWidget {
  final int value;
  final ValueChanged<int> onChanged;
  final List<RateLimitPreset> presets;
  final List<String> units;

  const RateLimit({
    super.key,
    required this.value,
    required this.onChanged,
    required this.presets,
    this.units = _defaultUnits,
  });

  @override
  State<RateLimit> createState() => _RateLimitState(value: value);
}

class _RateLimitState extends State<RateLimit> {
  int _value;
  String _unit;
  int _generation = 0;
  bool _expanded = false;

  _RateLimitState({required int value, String unit = 'sec'}) : _value = value, _unit = unit;

  void _selectUnit(String unit) {
    setState(() {
      _unit = unit;
      _expanded = false;
      _generation++;
    });
  }

  void _selectPreset(RateLimitPreset preset) {
    setState(() {
      _unit = preset.unit;
      _value = preset.value;
      _generation++;
    });
    widget.onChanged(preset.value);
  }

  void _onTextChanged(String text) {
    final n = int.tryParse(text);
    if (n == null) return;
    setState(() {
      _value = n;
    });
    widget.onChanged(n);
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
          initialValue: _value > 0 ? _value.toString() : '',
          keyboardType: TextInputType.number,
          inputFormatters: [FilteringTextInputFormatter.digitsOnly],
          decoration: InputDecoration(
            isDense: true,
            hintText: "0",
            helperText: "0 means no limit",
            suffixText: _unit,
            suffixIcon: IconButton(
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
        if (_expanded)
          Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            spacing: defaults.spacing,
            children: [
              if (widget.units.length > 1)
                LayoutBuilder(
                  builder: (context, constraints) {
                    final buttonWidth = (constraints.maxWidth - defaults.spacing * 2) / widget.units.length;
                    return Wrap(
                      spacing: defaults.spacing,
                      runSpacing: defaults.spacing,
                      alignment: WrapAlignment.start,
                      children:
                          widget.units
                              .map(
                                (u) => SizedBox(
                                  width: buttonWidth,
                                  child: OutlinedButton(
                                    onPressed: () => _selectUnit(u),
                                    style: OutlinedButton.styleFrom(
                                      padding: const EdgeInsets.symmetric(
                                        horizontal: 16,
                                        vertical: 8,
                                      ),
                                      shape: RoundedRectangleBorder(
                                        borderRadius: BorderRadius.circular(8),
                                      ),
                                    ),
                                    child: Text(u),
                                  ),
                                ),
                              )
                              .toList(),
                    );
                  },
                ),
              if (widget.presets.isNotEmpty) ...[
                Container(
                  width: double.infinity,
                  decoration: BoxDecoration(
                    border: Border(
                      bottom: BorderSide(color: Theme.of(context).dividerColor),
                    ),
                  ),
                  child: const Text('presets'),
                ),
                LayoutBuilder(
                  builder: (context, constraints) {
                    final buttonWidth = (constraints.maxWidth - defaults.spacing * 2) / 3.0;
                    return Wrap(
                      spacing: defaults.spacing,
                      runSpacing: defaults.spacing,
                      alignment: WrapAlignment.start,
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
            ],
          ),
      ],
    );
  }
}
