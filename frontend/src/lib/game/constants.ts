/**
 * constants.ts - ゲーム定数定義
 *
 * 場所、質問、タイマー時間などの固定値
 */

/**
 * S2, S6で使用する5つの場所
 */
export interface Place {
  id: string;
  name: string;
  description: string;
  imageUrl?: string;
}

export const FAVORITE_PLACES: Place[] = [
  {
    id: 'spiral_staircase',
    name: '螺旋階段を見上げる高い天井',
    description: '美しい螺旋階段が延びる、開放的な空間',
    imageUrl: '/images/places/spiral_staircase.jpg',
  },
  {
    id: 'main_hall_fireplace',
    name: 'メインホールの暖炉のレンガ',
    description: '歴史を感じさせる重厚な暖炉',
    imageUrl: '/images/places/fireplace.jpg',
  },
  {
    id: 'back_entrance_hinge',
    name: '裏玄関の扉の蝶番',
    description: '時を経た金属の輝き',
    imageUrl: '/images/places/hinge.jpg',
  },
  {
    id: 'front_entrance_door',
    name: '入口エントランスの扉',
    description: '多くの人を迎え入れてきた扉',
    imageUrl: '/images/places/entrance.jpg',
  },
  {
    id: 'upstairs_piano',
    name: '階上応接室のピアノ',
    description: '静かに佇む古いピアノ',
    imageUrl: '/images/places/piano.jpg',
  },
];

/**
 * S4診査室での10問の質問
 */
export interface Question {
  id: string;
  text: string;
  placeholder?: string;
  required: boolean;
  validation?: {
    minLength?: number;
    maxLength?: number;
    pattern?: RegExp;
  };
}

export const INTROSPECTION_QUESTIONS: Question[] = [
  {
    id: 'q1',
    text: '小学生の時、何に夢中でしたか？',
    placeholder: '例：サッカー、読書、絵を描くこと...',
    required: true,
    validation: {
      minLength: 2,
      maxLength: 100,
    },
  },
  {
    id: 'q2',
    text: 'その頃、誰を尊敬していましたか？',
    placeholder: '例：父、担任の先生、スポーツ選手...',
    required: true,
    validation: {
      minLength: 2,
      maxLength: 100,
    },
  },
  {
    id: 'q3',
    text: '中学生の時、どんな大人になりたいと思っていましたか？',
    placeholder: '例：医者になりたかった、自由に生きたかった...',
    required: true,
    validation: {
      minLength: 2,
      maxLength: 100,
    },
  },
  {
    id: 'q4',
    text: '高校生の時、一番大切だと思っていたものは何ですか？',
    placeholder: '例：友情、自由、夢...',
    required: true,
    validation: {
      minLength: 2,
      maxLength: 100,
    },
  },
  {
    id: 'q5',
    text: '今まで生きてきた中で、一番幸せだった瞬間はいつですか？',
    placeholder: '例：家族旅行、卒業式、初めての仕事...',
    required: true,
    validation: {
      minLength: 2,
      maxLength: 100,
    },
  },
  {
    id: 'q6',
    text: 'あなたが人生で一番後悔していることは何ですか？',
    placeholder: '例：やりたいことに挑戦しなかった、大切な人に伝えられなかった...',
    required: true,
    validation: {
      minLength: 2,
      maxLength: 100,
    },
  },
  {
    id: 'q7',
    text: '今、何のために生きていますか？',
    placeholder: '例：家族のため、夢を叶えるため、ただ生きるため...',
    required: true,
    validation: {
      minLength: 2,
      maxLength: 100,
    },
  },
  {
    id: 'q8',
    text: 'もし明日死ぬとしたら、今日何をしますか？',
    placeholder: '例：家族に会いに行く、やり残したことをする...',
    required: true,
    validation: {
      minLength: 2,
      maxLength: 100,
    },
  },
  {
    id: 'q9_life_goal',
    text: '人生の最期に達成したいことは何ですか？（これは重要な質問です）',
    placeholder: '例：世界を旅する、本を出版する、誰かの役に立つ...',
    required: true,
    validation: {
      minLength: 2,
      maxLength: 100,
    },
  },
  {
    id: 'q10_name',
    text: 'あなたの名前を教えてください',
    placeholder: 'あなたの名前',
    required: true,
    validation: {
      minLength: 1,
      maxLength: 50,
    },
  },
];

/**
 * タイマー時間（秒）
 */
export const TIMER_DURATIONS = {
  /** S6: 存在証明書探索 */
  S6_EXPLORATION: 7 * 60, // 7分
  /** S8: メインホール探索 */
  S8_EXPLORATION: 3 * 60, // 3分
};

/**
 * タイマー警告閾値（秒）
 */
export const TIMER_WARNING_THRESHOLDS = {
  S6: 60, // 残り1分
  S8: 30, // 残り30秒
};

/**
 * 画像類似度判定の閾値
 */
export const IMAGE_SIMILARITY_THRESHOLD = 0.5;

/**
 * メッセージ最大文字数
 */
export const MAX_MESSAGE_LENGTH = 120;

/**
 * VOICEVOX スピーカーID
 */
export const VOICEVOX_SPEAKER_IDS = {
  /** 青山龍星(しっとり) */
  担当者: 11,
};

/**
 * シーン名（表示用）
 */
export const SCENE_NAMES: Record<string, string> = {
  s0: '起動',
  s1: '1942年 - 初回訪問',
  s2: 'お気に入りの場所',
  s3: '移動指示',
  s4: '診査室',
  s5: '死亡届受理',
  s6: '存在証明書探索',
  s7: '2002年 - メッセージ',
  s8: 'メインホール',
  s9: '2025年 - あなたのメッセージ',
  gameover: 'ゲームオーバー',
};

/**
 * 無効な回答パターン（バリデーション用）
 */
export const INVALID_ANSWER_PATTERNS = [
  /^なし$/i,
  /^特になし$/i,
  /^ない$/i,
  /^わからない$/i,
  /^覚えていない$/i,
  /^無$/i,
  /^ー+$/,
  /^-+$/,
  /^\s*$/,
];

/**
 * 回答が無効かチェック
 */
export function isInvalidAnswer(answer: string): boolean {
  const trimmed = answer.trim();

  if (trimmed.length < 2) {
    return true;
  }

  return INVALID_ANSWER_PATTERNS.some((pattern) => pattern.test(trimmed));
}
