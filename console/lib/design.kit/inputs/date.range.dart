import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/timex.dart' as timex;
import 'package:retrovibed/design.kit/typography.dart' as typography;
import '../flutterx.dart';

class DateRangeInput extends StatefulWidget {
  final timex.Range value;
  final ValueChanged<timex.Range> onChanged;
  final DateTime firstDate;
  final DateTime lastDate;
  final bool autofocus;

  DateRangeInput({
    super.key,
    required this.value,
    required this.onChanged,
    DateTime? firstDate,
    DateTime? lastDate,
    this.autofocus = true,
  }) : firstDate = firstDate ?? timex.neginf,
       lastDate = lastDate ?? timex.inf;

  @override
  State<DateRangeInput> createState() => _DateRangeInputState(value);
}

class _DateRangeInputState extends State<DateRangeInput> {
  Widget _picker = ds.Empty;
  DateTime _current = timex.epoch;
  timex.Range _pending;

  // Owned explicitly (rather than relying on Focus(autofocus:)) because
  // _picker is swapped for a differently-keyed CalendarDatePicker every time
  // _showBegin/_showEnd runs — autofocus only fires once, on the wrapping
  // Focus widget's first attach, so it would never re-claim focus for a
  // picker shown by a later tap on the begin/end buttons (or by the
  // begin-pick auto-advancing to the end picker).
  final FocusNode _node = FocusNode(debugLabel: 'DateRangeInput picker');

  // parentScope (rather than the default closedLoop) lets Tab/Shift+Tab
  // escape to whatever's next outside this widget — e.g. a field's preset
  // suggestion list shown alongside the picker — instead of endlessly
  // cycling begin/end/calendar in place.
  final FocusScopeNode _focusScope = FocusScopeNode(
    traversalEdgeBehavior: TraversalEdgeBehavior.parentScope,
  );

  _DateRangeInputState(this._pending);

  @override
  void initState() {
    super.initState();
    postframe(() => _showBegin(focus: widget.autofocus));
  }

  @override
  void dispose() {
    _focusScope.dispose();
    _node.dispose();
    super.dispose();
  }

  void _showBegin({bool focus = true}) {
    final firstDate = timex.min([widget.firstDate, _pending.begin]);
    final lastDate = timex.max([widget.lastDate, _pending.end]);

    setState(() {
      _current = _pending.begin;
      _picker = CalendarDatePicker(
        key: const ValueKey('begin'),
        initialDate: timex.max([_pending.begin.toLocal(), timex.now()]),
        firstDate: firstDate.toLocal(),
        lastDate: lastDate.toLocal(),
        onDateChanged: (d) {
          setState(() {
            _pending = timex.Range(d.toUtc(), _pending.end);
          });
          _showEnd();
        },
      );
    });
    if (focus) postframe(() => _node.requestFocus());
  }

  void _showEnd({bool focus = true}) {
    final firstDate = timex.min([widget.firstDate, _pending.begin]);
    final lastDate = timex.max([widget.lastDate, _pending.end]);
    setState(() {
      _current = _pending.end;
      _picker = CalendarDatePicker(
        key: const ValueKey('end'),
        initialDate: _pending.end.toLocal(),
        firstDate: firstDate.toLocal(),
        lastDate: lastDate.toLocal(),
        onDateChanged: (d) {
          final next = timex.Range(_pending.begin, d.toUtc());
          setState(() => _pending = next);
          widget.onChanged(next);
        },
      );
    });
    if (focus) postframe(() => _node.requestFocus());
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final theme = Theme.of(context);
    final activestyle = TextButton.styleFrom(
      backgroundColor: theme.colorScheme.primaryContainer,
    );

    return FocusScope(
      node: _focusScope,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            spacing: defaults.spacing,
            children: [
              Expanded(
                child: TextButton(
                  onPressed: _showBegin,
                  style: timex.neginf != _current && _current == _pending.begin ? activestyle : null,
                  child: typography.Timestamp(_pending.begin, neginf: Text("")),
                ),
              ),
              const Text('–'),
              Expanded(
                child: TextButton(
                  onPressed: _showEnd,
                  style: timex.inf != _current && _current == _pending.end ? activestyle : null,
                  child: typography.Timestamp(_pending.end),
                ),
              ),
            ],
          ),
          Focus(focusNode: _node, autofocus: widget.autofocus, child: _picker),
        ],
      ),
    );
  }
}
