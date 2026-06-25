import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './field.dart';
import './suggestion.list.dart';
import 'parser.results.dart';

sealed class ParserState extends StatelessWidget {
  const ParserState({super.key});
  (ParserState, num) consume(TextEditingController ctrl);
}

class Context {
  final int offset;
  final int lastOffset;
  final String partial;
  final void Function(
    Context ctx,
    TextRange range,
    String contents, {
    ParserResult? completed,
  })
  replacement;
  final List<Field<dynamic>> fields;
  final GlobalKey<SuggestionListState> suggestionKey;
  const Context(
    this.fields,
    this.offset,
    this.partial,
    this.replacement,
    this.suggestionKey, {
    this.lastOffset = 0,
  });

  num remaining(TextEditingController ctrl) => ctrl.selection.baseOffset - lastOffset;
}

class Query extends ParserState {
  final Context ctx;
  Query(this.ctx);

  @override
  (ParserState, num) consume(TextEditingController ctrl) {
    final cursor = ctrl.selection.baseOffset;
    if (cursor < ctx.lastOffset) {
      final resetTo = ctx.offset < 0 ? 0 : ctx.offset;
      final next = Context(
        ctx.fields,
        ctx.offset,
        '',
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: resetTo,
      );
      return (Query(next), next.remaining(ctrl));
    }

    if (cursor == ctx.lastOffset) {
      final startPos = ctx.lastOffset - ctx.partial.length;
      final actual = cursor > 0 ? ctrl.text.substring(startPos, cursor) : '';
      if (actual == ctx.partial) {
        return (this, 0);
      }

      final next = Context(
        ctx.fields,
        ctx.offset,
        '',
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: startPos,
      );
      return (Query(next), cursor - startPos);
    }

    final char = ctrl.text[ctx.lastOffset];
    final nextLastOffset = ctx.lastOffset + 1;

    if (char == '@') {
      final next = Context(
        ctx.fields,
        ctx.lastOffset,
        '',
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: nextLastOffset,
      );
      return (Input(next), cursor - nextLastOffset);
    }

    final next = Context(
      ctx.fields,
      ctx.offset,
      ctx.partial + char,
      ctx.replacement,
      ctx.suggestionKey,
      lastOffset: nextLastOffset,
    );
    return (Query(next), cursor - nextLastOffset);
  }

  @override
  Widget build(BuildContext context) => ds.Empty;
}

class Input extends ParserState {
  final Context ctx;
  Input(this.ctx);

  @override
  (ParserState, num) consume(TextEditingController ctrl) {
    final cursor = ctrl.selection.baseOffset;
    if (cursor <= ctx.lastOffset) {
      final next = Context(
        ctx.fields,
        ctx.offset,
        '',
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: ctx.offset,
      );
      return (Query(next), next.remaining(ctrl));
    }

    final char = ctrl.text[ctx.lastOffset];
    final nextLastOffset = ctx.lastOffset + 1;

    if (char == ':') {
      final substring = ctrl.text.substring(ctx.offset + 1, ctx.lastOffset);
      final field = ctx.fields.where((f) => f.name == substring).firstOrNull;
      if (field == null) {
        final next = Context(
          ctx.fields,
          ctx.offset,
          substring,
          ctx.replacement,
          ctx.suggestionKey,
          lastOffset: nextLastOffset,
        );
        return (UnknownFieldError(next), cursor - nextLastOffset);
      }

      final defaultValue = field.autocomplete;
      final upd = ctx.fields
          .map(
            (f) => f.name == field.name ? field.withCurrent(defaultValue ?? field.defaultValue) : f,
          )
          .toList();

      if (defaultValue != null) {
        final completed = field.of(defaultValue);
        ctx.replacement(
          ctx,
          TextRange(start: ctx.offset, end: nextLastOffset),
          '',
          completed: completed,
        );

        final next = Context(
          upd,
          ctx.offset,
          '',
          ctx.replacement,
          ctx.suggestionKey,
          lastOffset: nextLastOffset,
        );
        return (Query(next), 0);
      }

      final next = Context(
        upd,
        ctx.offset,
        "@${substring}:",
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: nextLastOffset,
      );
      return (Value(next, field), cursor - nextLastOffset);
    }

    final next = Context(
      ctx.fields,
      ctx.offset,
      ctx.partial + char,
      ctx.replacement,
      ctx.suggestionKey,
      lastOffset: nextLastOffset,
    );
    return (Input(next), cursor - nextLastOffset);
  }

  @override
  Widget build(BuildContext context) {
    final matches = ctx.fields.where((f) => f.available && f.name.startsWith(ctx.partial)).toList();
    if (matches.isEmpty) return ds.Empty;

    final key = ctx.suggestionKey;
    return SuggestionList([
      for (final m in matches)
        (
          Text(m.name),
          () => ctx.replacement(
            ctx,
            TextRange(
              start: ctx.offset,
              end: ctx.offset + ctx.partial.length + 1,
            ),
            '@${m.name}:',
          ),
        ),
    ], key: key);
  }
}

class UnknownFieldError extends ParserState {
  final Context ctx;
  UnknownFieldError(this.ctx);

  @override
  (ParserState, num) consume(TextEditingController ctrl) {
    final cursor = ctrl.selection.baseOffset;
    if (cursor <= 0) {
      final next = Context(
        ctx.fields,
        -1,
        '',
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: cursor,
      );
      return (Query(next), next.remaining(ctrl));
    }

    final text = ctrl.text;

    // If the @ anchor is gone, bail out to Query.
    if (ctx.offset < 0 || ctx.offset >= text.length || text[ctx.offset] != '@') {
      final next = Context(
        ctx.fields,
        cursor - 1,
        '',
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: cursor,
      );
      return (Query(next), next.remaining(ctrl));
    }

    // Re-read the segment between @ and the current cursor.
    final segment = text.substring(
      ctx.offset + 1,
      cursor.clamp(ctx.offset + 1, text.length),
    );
    final colonIdx = segment.indexOf(':');

    if (colonIdx == -1) {
      final next = Context(
        ctx.fields,
        ctx.offset,
        segment,
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: cursor,
      );
      return (Input(next), next.remaining(ctrl));
    }

    final partial = segment.substring(0, colonIdx);
    final field = ctx.fields.where((f) => f.name == partial).firstOrNull;
    if (field == null) {
      final next = Context(
        ctx.fields,
        ctx.offset,
        partial,
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: cursor,
      );
      return (UnknownFieldError(next), next.remaining(ctrl));
    }

    final next = Context(
      ctx.fields,
      ctx.offset,
      segment.substring(colonIdx + 1),
      ctx.replacement,
      ctx.suggestionKey,
      lastOffset: cursor,
    );
    return (Value(next, field), next.remaining(ctrl));
  }

  @override
  Widget build(BuildContext context) => ds.Empty;

  @override
  String toString({DiagnosticLevel minLevel = DiagnosticLevel.info}) =>
      'UnknownFieldError: unknown field "${ctx.partial}"';
}

class Value extends ParserState {
  final Context ctx;
  final Field<dynamic> field;

  Value(this.ctx, this.field);

  @override
  (ParserState, num) consume(TextEditingController ctrl) {
    final cursor = ctrl.selection.baseOffset;
    if (cursor <= 0 || cursor <= ctx.offset) {
      final next = Context(
        ctx.fields,
        -1,
        '',
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: cursor < 0 ? 0 : cursor,
      );
      return (Query(next), next.remaining(ctrl));
    }

    final segment = ctrl.text.substring(
      ctx.offset,
      cursor.clamp(ctx.offset, ctrl.text.length),
    );
    final colonIdx = segment.indexOf(':');

    if (colonIdx == -1) {
      // Colon was deleted — go back to Input with the field name partial.
      final fieldPartial = segment.length > 1 ? segment.substring(1) : '';
      final next = Context(
        ctx.fields,
        ctx.offset,
        fieldPartial,
        ctx.replacement,
        ctx.suggestionKey,
        lastOffset: cursor,
      );
      return (Input(next), next.remaining(ctrl));
    }

    final next = Context(
      ctx.fields,
      ctx.offset,
      segment,
      ctx.replacement,
      ctx.suggestionKey,
      lastOffset: cursor,
    );
    return (Value(next, field), next.remaining(ctrl));
  }

  String get _valuePartial {
    final colon = ctx.partial.indexOf(':');
    return colon == -1 ? ctx.partial : ctx.partial.substring(colon + 1);
  }

  void _commit(dynamic v) {
    ctx.replacement(
      ctx,
      TextRange(start: ctx.offset, end: ctx.offset + ctx.partial.length),
      "",
      completed: field.of(v) as ParserResult?,
    );
  }

  @override
  Widget build(BuildContext context) {
    final valuePartial = _valuePartial;
    final suggestions = field.suggestions(valuePartial);
    final theme = Theme.of(context);
    final renderWidget = field.render(field.parse(valuePartial), _commit);
    final hasSuggestions = suggestions.isNotEmpty;
    final hasRender = renderWidget != ds.Empty;

    if (!hasSuggestions && !hasRender) return ds.Empty;

    if (!hasSuggestions) return renderWidget;

    final suggestionItems = [
      for (final s in suggestions)
        (
          Text(s.label) as Widget,
          () => _commit(
            field.parse(s.completion.substring(s.completion.indexOf(':') + 1)),
          ),
        ),
    ];

    final key = ctx.suggestionKey;

    if (!hasRender) {
      return SuggestionList(suggestionItems, key: key);
    }

    return Material(
      elevation: 4,
      color: theme.colorScheme.surface,
      borderRadius: BorderRadius.circular(4),
      clipBehavior: Clip.antiAlias,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          renderWidget,
          SuggestionList(suggestionItems, key: key),
        ],
      ),
    );
  }
}

class Parser {
  ParserState current;

  Parser(
    List<Field<dynamic>> fields,
    void Function(
      Context ctx,
      TextRange range,
      String contents, {
      ParserResult? completed,
    })
    replacement,
    GlobalKey<SuggestionListState> suggestionKey,
  ) : current = Query(Context(fields, -1, '', replacement, suggestionKey));

  /// Restores a single field in the current parser context, making it
  /// available in suggestions again (e.g. after its chip is deleted).
  void replace(Field<dynamic> field) {
    final ctx = switch (current) {
      Query q => q.ctx,
      Input i => i.ctx,
      Value v => v.ctx,
      _ => null,
    };
    if (ctx == null) return;
    final fields = ctx.fields.map((f) => f.name == field.name ? field : f).toList();
    final next = Context(
      fields,
      ctx.offset,
      ctx.partial,
      ctx.replacement,
      ctx.suggestionKey,
      lastOffset: ctx.lastOffset,
    );
    current = switch (current) {
      Query _ => Query(next),
      Input _ => Input(next),
      Value v => Value(next, v.field),
      _ => current,
    };
  }

  ParserState consume(TextEditingController ctrl) {
    if (ctrl.selection.baseOffset < 0) return current;
    var next = current;
    num remaining;
    do {
      (next, remaining) = next.consume(ctrl);
    } while (remaining > 0);
    return current = next;
  }
}
