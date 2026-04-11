import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/lucene.dart' as lucene;

void main() {
  group('parse — bare words', () {
    test('single word', () {
      final q = lucene.parse('hello');
      expect(q.clauses, hasLength(1));
      final t = q.clauses.first.node as lucene.Term;
      expect(t.field, isEmpty);
      expect(t.value, 'hello');
    });

    test('multiple words produce multiple clauses', () {
      final q = lucene.parse('big buck bunny');
      expect(q.clauses, hasLength(3));
      expect((q.clauses[0].node as lucene.Term).value, 'big');
      expect((q.clauses[1].node as lucene.Term).value, 'buck');
      expect((q.clauses[2].node as lucene.Term).value, 'bunny');
    });

    test('bare word defaults to should', () {
      final q = lucene.parse('hello');
      expect(q.clauses.first.occur, lucene.Occur.should);
    });
  });

  group('parse — phrases', () {
    test('quoted phrase is single term', () {
      final q = lucene.parse('"big buck bunny"');
      expect(q.clauses, hasLength(1));
      final t = q.clauses.first.node as lucene.Term;
      expect(t.field, isEmpty);
      expect(t.value, 'big buck bunny');
    });

    test('phrase with escape', () {
      final q = lucene.parse(r'"say \"hello\""');
      final t = q.clauses.first.node as lucene.Term;
      expect(t.value, 'say "hello"');
    });
  });

  group('parse — field terms', () {
    test('field:value', () {
      final q = lucene.parse('type:movie');
      final t = q.clauses.first.node as lucene.Term;
      expect(t.field, 'type');
      expect(t.value, 'movie');
    });

    test('field with no value produces empty term', () {
      final q = lucene.parse('hd:');
      final t = q.clauses.first.node as lucene.Term;
      expect(t.field, 'hd');
      expect(t.value, '');
    });
  });

  group('parse — ranges', () {
    test('inclusive range [x TO y]', () {
      final q = lucene.parse('date:[2025-01-01 TO 2026-01-01]');
      final r = q.clauses.first.node as lucene.Range;
      expect(r.field, 'date');
      expect(r.min, '2025-01-01');
      expect(r.max, '2026-01-01');
      expect(r.minInclusive, isTrue);
      expect(r.maxInclusive, isTrue);
    });

    test('exclusive range {x TO y}', () {
      final q = lucene.parse('date:{2025-01-01 TO 2026-01-01}');
      final r = q.clauses.first.node as lucene.Range;
      expect(r.minInclusive, isFalse);
      expect(r.maxInclusive, isFalse);
    });

    test('open-ended range [* TO y]', () {
      final q = lucene.parse('size:[* TO 1GiB]');
      final r = q.clauses.first.node as lucene.Range;
      expect(r.min, isNull);
      expect(r.max, '1GiB');
    });

    test('comparison > becomes Range', () {
      final q = lucene.parse('peers:>5');
      final r = q.clauses.first.node as lucene.Range;
      expect(r.field, 'peers');
      expect(r.min, '5');
      expect(r.minInclusive, isFalse);
      expect(r.max, isNull);
    });

    test('comparison >= becomes inclusive Range', () {
      final q = lucene.parse('peers:>=5');
      final r = q.clauses.first.node as lucene.Range;
      expect(r.minInclusive, isTrue);
    });

    test('comparison < becomes Range on max', () {
      final q = lucene.parse('size:<500MiB');
      final r = q.clauses.first.node as lucene.Range;
      expect(r.min, isNull);
      expect(r.max, '500MiB');
      expect(r.maxInclusive, isFalse);
    });

    test('ISO-8601 timestamps with colons parse as complete bounds', () {
      final q = lucene.parse(
        'date:[2020-03-01T00:00:00.000Z TO 2021-03-01T00:00:00.000Z]',
      );
      final r = q.clauses.first.node as lucene.Range;
      expect(r.field, 'date');
      expect(r.min, '2020-03-01T00:00:00.000Z');
      expect(r.max, '2021-03-01T00:00:00.000Z');
    });

    test('ISO-8601 open-ended range [ts TO *] preserves min', () {
      final q = lucene.parse('date:[2020-03-01T00:00:00.000Z TO *]');
      final r = q.clauses.first.node as lucene.Range;
      expect(r.min, '2020-03-01T00:00:00.000Z');
      expect(r.max, isNull);
    });
  });

  group('parse — occurrence modifiers', () {
    test('+ prefix = must', () {
      final q = lucene.parse('+hd:true');
      expect(q.clauses.first.occur, lucene.Occur.must);
    });

    test('- prefix = mustNot', () {
      final q = lucene.parse('-hd:true');
      expect(q.clauses.first.occur, lucene.Occur.mustNot);
    });

    test('NOT keyword = mustNot', () {
      final q = lucene.parse('NOT hd:true');
      expect(q.clauses.first.occur, lucene.Occur.mustNot);
    });
  });

  group('parse — groups', () {
    test('parenthesised sub-query becomes Group', () {
      final q = lucene.parse('(type:movie OR type:series)');
      final g = q.clauses.first.node as lucene.Group;
      expect(g.query.clauses, hasLength(2));
    });
  });

  group('parse — AND/OR connectors are skipped', () {
    test('AND skipped as connector', () {
      final q = lucene.parse('type:movie AND hd:true');
      expect(q.clauses, hasLength(2));
    });

    test('OR skipped as connector', () {
      final q = lucene.parse('type:movie OR type:series');
      expect(q.clauses, hasLength(2));
    });
  });

  group('parse — mixed query', () {
    test('field filters + free text', () {
      final q = lucene.parse('type:movie -hd:true "buck bunny"');
      expect(q.clauses, hasLength(3));
      final term = q.clauses[0].node as lucene.Term;
      final hd = q.clauses[1];
      final phrase = q.clauses[2].node as lucene.Term;
      expect(term.field, 'type');
      expect(hd.occur, lucene.Occur.mustNot);
      expect(phrase.value, 'buck bunny');
    });
  });

  group('parse — empty input', () {
    test('empty string produces empty query', () {
      final q = lucene.parse('');
      expect(q.clauses, isEmpty);
    });

    test('only whitespace produces empty query', () {
      final q = lucene.parse('   ');
      expect(q.clauses, isEmpty);
    });
  });
}
