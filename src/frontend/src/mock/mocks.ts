import { IMovie } from '../models/MovieInfo';
import image1 from '../images/image_1.jpeg';
import image2 from '../images/image_2.jpeg';
import image3 from '../images/image_3.jpeg';
import image4 from '../images/image_4.jpeg';
import video1 from '../images/video01.webm';
import video2 from '../images/video02.webm';
import video3 from '../images/video03.webm';

export const cards: IMovie[] = [
  {
    id: '0',
    name: 'Capitan Marvel',
    img:
      'https://upload.wikimedia.org/wikipedia/ru/0/07/Captain_Marvel_film_logo.jpg',
    src: video1,
    info: {
      avalibility: 0.3,
      year: 2000,
      genres: ['Action', 'Sci-Fi'],
      country: ['USA'],
      rating: 3.4,
      views: 123000,
      length: 123,
      pgRating: 'PG-13',
      description:
        "Carol Danvers becomes one of the universe's most powerful heroes when Earth is caught in the middle of a galactic war between two alien races.",
      photos: [image1, image2, image3, image4],
      videos: [video1, video2, video3],
      comments: [
        {
          id: '1',
          text: 'hello',
          username: 'Alex',
          avatar: 'https://randomuser.me/api/portraits/men/2.jpg',
        },
        {
          id: '2',
          text: 'nice movie',
          username: 'Lee',
          avatar: 'https://randomuser.me/api/portraits/women/2.jpg',
        },
      ],
    },
  },
  {
    id: '1',
    name: 'Hobbit',
    img:
      'https://images-na.ssl-images-amazon.com/images/I/7145Wo9GjlL._AC_SL1006_.jpg',
    src: video2,
    info: {
      avalibility: 0.1,
      year: 2004,
      genres: ['Fantasy', 'Adventure'],
      rating: 4.9,
      views: 9123000,
      length: 123,
      pgRating: 'PG-6',
      description:
        'A reluctant Hobbit, Bilbo Baggins, sets out to the Lonely Mountain with a spirited group of dwarves to reclaim their mountain home, and the gold within it from the dragon Smaug. A reluctant Hobbit, Bilbo Baggins, sets out to the Lonely Mountain with a spirited group of dwarves to reclaim their mountain home, and the gold within it from the dragon Smaug',
    },
  },
  {
    id: '2',
    name: 'Star Trek',
    img: 'https://m.media-amazon.com/images/I/615RWFNlXDL._AC_SY741_.jpg',
    src: video1,
    info: {
      avalibility: 0.5,
      year: 1994,
      genres: ['Action', 'Sci-Fi'],
      rating: 1.4,
      views: 987654,
      length: 123,
      pgRating: 'PG-13',
    },
  },
  {
    id: '3',
    name: 'Capitan Marvel',
    img:
      'https://upload.wikimedia.org/wikipedia/ru/0/07/Captain_Marvel_film_logo.jpg',
    src: video1,
    info: {
      avalibility: 0.8,
      year: 2000,
      genres: ['Action', 'Sci-Fi'],
      rating: 3.4,
      views: 123000,
      length: 123,
      pgRating: 'PG-13',
    },
  },
  {
    id: '4',
    name: 'Hobbit',
    img:
      'https://images-na.ssl-images-amazon.com/images/I/7145Wo9GjlL._AC_SL1006_.jpg',
    src: video1,
    info: {
      avalibility: 0.23,
      year: 2004,
      genres: ['Action', 'Sci-Fi'],
      rating: 4.9,
      views: 9123000,
      length: 123,
      pgRating: 'PG-6',
    },
  },
  {
    id: '5',
    name: 'Star Trek',
    img: 'https://m.media-amazon.com/images/I/615RWFNlXDL._AC_SY741_.jpg',
    src: video1,
    info: {
      avalibility: 0.56,
      year: 1994,
      genres: ['Action', 'Sci-Fi'],
      rating: 1.4,
      views: 987654,
      length: 123,
      pgRating: 'PG-13',
    },
  },
  {
    id: '6',
    name: 'Capitan Marvel',
    img:
      'https://upload.wikimedia.org/wikipedia/ru/0/07/Captain_Marvel_film_logo.jpg',
    src: video1,
    info: {
      avalibility: 0.51,
      year: 2000,
      genres: ['Action', 'Sci-Fi'],
      rating: 3.4,
      views: 123000,
      length: 123,
      pgRating: 'PG-13',
    },
  },
  {
    id: '7',
    name: 'Hobbit',
    img:
      'https://images-na.ssl-images-amazon.com/images/I/7145Wo9GjlL._AC_SL1006_.jpg',
    src: video1,
    info: {
      avalibility: 0.11,
      year: 2004,
      genres: ['Action', 'Sci-Fi'],
      rating: 4.9,
      views: 9123000,
      length: 123,
      pgRating: 'PG-6',
    },
  },
  {
    id: '8',
    name: 'Star Trek',
    img: 'https://m.media-amazon.com/images/I/615RWFNlXDL._AC_SY741_.jpg',
    src: video1,
    info: {
      avalibility: 0.71,
      year: 1994,
      genres: ['Action', 'Sci-Fi'],
      rating: 1.4,
      views: 987654,
      length: 123,
      pgRating: 'PG-13',
    },
  },
  {
    id: '9',
    name: 'Capitan Marvel',
    img:
      'https://upload.wikimedia.org/wikipedia/ru/0/07/Captain_Marvel_film_logo.jpg',
    src: video1,
    info: {
      avalibility: 0.1,
      year: 2000,
      genres: ['Action', 'Sci-Fi'],
      rating: 3.4,
      views: 123000,
      length: 123,
      pgRating: 'PG-13',
    },
  },
  {
    id: '10',
    name: 'Hobbit',
    img:
      'https://images-na.ssl-images-amazon.com/images/I/7145Wo9GjlL._AC_SL1006_.jpg',
    src: video1,
    info: {
      avalibility: 0.01,
      year: 2004,
      genres: ['Western', 'Drama'],
      rating: 4.9,
      views: 9123000,
      length: 123,
      pgRating: 'PG-6',
    },
  },
  {
    id: '11',
    name: 'Star Trek',
    img: 'https://m.media-amazon.com/images/I/615RWFNlXDL._AC_SY741_.jpg',
    src: video1,
    info: {
      avalibility: 0.21,
      year: 1994,
      genres: ['Detective'],
      rating: 1.4,
      views: 987654,
      length: 123,
      pgRating: 'PG-13',
    },
  },
];
